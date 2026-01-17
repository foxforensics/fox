package loader

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/cuhsat/go-mmap"
	"github.com/sourcegraph/conc"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heap"
	"github.com/cuhsat/fox/v4/internal/pkg/types/register"
)

const stdin = "-"

type Options struct {
	Limit    *types.Limits
	Filter   *types.Filters
	File     string
	Input    string
	Password string
	Profile  int
	Verbose  int
}

type Loader struct {
	sync.RWMutex
	opts  *Options
	size  int64
	paths []string
	heaps chan *heap.Heap
}

func New(opts *Options) *Loader {
	return &Loader{
		opts:  opts,
		heaps: make(chan *heap.Heap, opts.Profile),
	}
}

func (ldr *Loader) Load(paths []string) <-chan *heap.Heap {
	if len(ldr.opts.File) > 0 {
		b, err := os.ReadFile(ldr.opts.File)

		if err != nil {
			log.Fatalln(err)
		}

		lines := strings.Split(strings.TrimSpace(string(b)), "\n")

		if ldr.opts.Verbose > 0 {
			for _, l := range lines {
				log.Printf("add path %s \n", l)
			}
		}

		paths = append(paths, lines...)
	}

	go func() {
		defer close(ldr.heaps)

		if isFilePiped(os.Stdin) {
			paths = append(paths, stdin)
		}

		if len(ldr.opts.Input) > 0 {
			ldr.createHeap("input", []byte(ldr.opts.Input))
		}

		for _, path := range paths {
			if path == stdin {
				if !isFilePiped(os.Stdin) {
					log.Fatalln("stdin not open")
				}

				buf, err := io.ReadAll(bufio.NewReader(os.Stdin))

				if err != nil {
					log.Fatalln(err)
				}

				ldr.createHeap(stdin, buf)
				break
			}

			path, part := data.SplitPart(path)

			_, err := os.Stat(path)

			if ldr.opts.Verbose > 0 && errors.Is(err, os.ErrNotExist) {
				log.Printf("looked for %s\n", path)
			}

			ldr.loadPath(path, part)
		}
	}()

	return ldr.heaps
}

func (ldr *Loader) Exit() {
	ldr.RLock()

	if ldr.opts.Verbose > 0 {
		log.Printf("size %s\n", text.Humanize(ldr.size))
	}

	ldr.RUnlock()
}

func (ldr *Loader) loadPath(path, part string) {
	match, err := doublestar.FilepathGlob(path)

	if err != nil {
		log.Println(err)
		return
	}

	var wg conc.WaitGroup

	for _, path := range match {
		fi, err := os.Stat(path)

		if err != nil {
			log.Println(err)
			continue
		}

		wg.Go(func() {
			if fi.IsDir() {
				ldr.loadDir(path, part)
			} else {
				ldr.loadFile(path, part)
			}
		})
	}

	wg.Wait()
}

func (ldr *Loader) loadDir(path, part string) {
	dir, err := os.ReadDir(path)

	if err != nil {
		log.Println(err)
		return
	}

	var wg conc.WaitGroup

	for _, f := range dir {
		if !f.IsDir() {
			wg.Go(func() {
				ldr.loadFile(filepath.Join(path, f.Name()), part)
			})
		}
	}

	wg.Wait()
}

func (ldr *Loader) loadFile(path, part string) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0x400)

	if err != nil {
		log.Println(err)
		return
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	fi, err := f.Stat()

	if err != nil {
		log.Println(err)
		return
	}

	// empty files will cause issues
	if fi.Size() == 0 {
		ldr.createHeap(path, []byte{})
		return
	}

	b, err := mmap.Map(f, mmap.RDONLY, 0)

	if err != nil {
		log.Println(err)
		return
	}

	if ldr.opts.Verbose > 2 {
		log.Printf("mapped file %s\n", path)
	}

	ldr.processData(path, part, b)
}

func (ldr *Loader) processData(path, part string, b []byte) {
	// 1. deflate data
	b = ldr.deflateData(b)

	// 2. extract data
	if ldr.extractData(path, part, b) {
		return
	}

	// 3a. ingest data
	b = ldr.ingestData(b)

	// 3b. convert data
	b = ldr.convertData(b)

	// filter for specific streams
	if strings.Contains(path, part) {
		ldr.createHeap(path, b)
	}
}

func (ldr *Loader) extractData(path, part string, b []byte) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Println("archive corrupt or password wrong")
			return
		}
	}()

	for _, a := range register.Archives {
		if a.Detect(b) {
			if ldr.opts.Verbose > 1 {
				log.Printf("archive detected %s\n", a.Name)
			}

			var wg conc.WaitGroup

			for _, e := range a.Extract(b, path, ldr.opts.Password) {
				wg.Go(func() {
					if ldr.opts.Verbose > 2 {
						log.Printf("stream detected %s\n", e.Path)
					}

					ldr.processData(e.Path, part, e.Data)
				})
			}

			wg.Wait()

			return true
		}
	}

	return false
}

func (ldr *Loader) deflateData(b []byte) []byte {
	for _, d := range register.Deflates {
		if d.Detect(b) {
			if ldr.opts.Verbose > 1 {
				log.Printf("deflate detected %s\n", d.Name)
			}

			r, err := d.Deflate(b)

			if err != nil {
				log.Println(err)
			}

			return r
		}
	}

	return b
}

func (ldr *Loader) convertData(b []byte) []byte {
	for _, c := range register.Converts {
		if c.Detect(b) {
			if ldr.opts.Verbose > 1 {
				log.Printf("convert detected %s\n", c.Name)
			}

			r, err := c.Convert(b)

			if err != nil {
				log.Println(err)
			}

			return r
		}
	}

	return b
}

func (ldr *Loader) ingestData(b []byte) []byte {
	for _, c := range register.Images {
		if c.Detect(b) {
			if ldr.opts.Verbose > 1 {
				log.Printf("image detected %s\n", c.Name)
			}

			r, err := c.Ingest(b)

			if err != nil {
				log.Println(err)
			}

			return r
		}
	}

	return b
}

func (ldr *Loader) createHeap(s string, b []byte) {
	ldr.Lock()
	defer ldr.Unlock()

	if slices.Contains(ldr.paths, s) {
		return // already loaded
	}

	ldr.size += int64(len(b))
	ldr.paths = append(ldr.paths, s)
	ldr.heaps <- heap.New(&heap.Context{
		Name:   s,
		Limit:  ldr.opts.Limit,
		Filter: ldr.opts.Filter,
	}, b)

	if ldr.opts.Verbose > 1 {
		log.Printf("loaded heap %s\n", s)
	}
}

func isFilePiped(f *os.File) bool {
	fi, err := f.Stat()

	if err != nil {
		log.Fatalln(err)
	}

	return (fi.Mode() & os.ModeCharDevice) != os.ModeCharDevice
}
