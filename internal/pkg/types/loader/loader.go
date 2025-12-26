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
	"github.com/cuhsat/fox/v4/internal/pkg/data"
	"github.com/edsrzf/mmap-go"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heap"
	"github.com/cuhsat/fox/v4/internal/pkg/types/register"
)

const stdin = "-"

type Options struct {
	Limit    *types.Limits
	Filter   *types.Filters
	Queue    uint
	Input    string
	Password string
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
		heaps: make(chan *heap.Heap, opts.Queue),
	}
}

func (ldr *Loader) Load(paths []string) <-chan *heap.Heap {
	go func() {
		defer close(ldr.heaps)

		if isFilePiped(os.Stdin) {
			paths = append(paths, stdin)
		}

		if len(ldr.opts.Input) > 0 {
			ldr.createFromInput([]byte(ldr.opts.Input))
		}

		for _, path := range paths {
			if path == stdin {
				ldr.createFromStdin()
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

	var wg sync.WaitGroup

	for _, path := range match {
		fi, err := os.Stat(path)

		if err != nil {
			log.Println(err)
			continue
		}

		wg.Add(1)

		go func() {
			if fi.IsDir() {
				ldr.loadDir(path, part)
			} else {
				ldr.loadFile(path, part)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func (ldr *Loader) loadDir(path, part string) {
	dir, err := os.ReadDir(path)

	if err != nil {
		log.Println(err)
		return
	}

	var wg sync.WaitGroup

	for _, f := range dir {
		if !f.IsDir() {
			wg.Add(1)

			go func() {
				ldr.loadFile(filepath.Join(path, f.Name()), part)
				wg.Done()
			}()
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
		ldr.createFromData(path, []byte{})
		return
	}

	b, err := mmap.Map(f, mmap.RDONLY, 0)

	if err != nil {
		log.Println(err)
		return
	}

	ldr.processData(path, part, b, false)
}

func (ldr *Loader) processData(path, part string, b []byte, data bool) {
	var ok bool

	// 1. deflate data
	if b, ok = ldr.deflateData(b); ok {
		data = true
	}

	// 2. extract data
	if ldr.extractData(path, part, b) {
		return
	}

	// 3. format data
	if b, ok = ldr.formatData(b); ok {
		data = true
	}

	// filter for specific streams
	if strings.Contains(path, part) {
		if data {
			ldr.createFromData(path, b)
		} else {
			ldr.createFromFile(path, b)
		}
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

			var wg sync.WaitGroup

			for _, e := range a.Extract(b, path, ldr.opts.Password) {
				wg.Add(1)

				go func() {
					if ldr.opts.Verbose > 2 {
						log.Printf("stream detected %s\n", e.Path)
					}

					ldr.processData(e.Path, part, e.Data, true)

					wg.Done()
				}()
			}

			wg.Wait()

			return true
		}
	}

	return false
}

func (ldr *Loader) deflateData(b []byte) ([]byte, bool) {
	for _, d := range register.Deflates {
		if d.Detect(b) {
			if ldr.opts.Verbose > 1 {
				log.Printf("deflate detected %s\n", d.Name)
			}

			r, err := d.Deflate(b)

			if err != nil {
				log.Println(err)
			}

			return r, true
		}
	}

	return b, false
}

func (ldr *Loader) formatData(b []byte) ([]byte, bool) {
	for _, c := range register.Converts {
		if c.Detect(b) {
			if ldr.opts.Verbose > 1 {
				log.Printf("convert detected %s\n", c.Name)
			}

			r, err := c.Convert(b)

			if err != nil {
				log.Println(err)
			}

			return r, true
		}
	}

	return b, false
}

func (ldr *Loader) createFromFile(path string, b []byte) {
	ldr.createHeap(path, types.Regular, b)

	if ldr.opts.Verbose > 1 {
		log.Printf("loaded heap from file %s\n", path)
	}
}

func (ldr *Loader) createFromData(name string, b []byte) {
	ldr.createHeap(name, types.Deflate, b)

	if ldr.opts.Verbose > 1 {
		log.Printf("loaded heap from data %s\n", name)
	}
}

func (ldr *Loader) createFromInput(b []byte) {
	ldr.createHeap("input", types.Defined, b)

	if ldr.opts.Verbose > 1 {
		log.Println("loaded heap from input")
	}
}

func (ldr *Loader) createFromStdin() {
	if !isFilePiped(os.Stdin) {
		log.Fatalln("stdin not open")
	}

	buf, err := io.ReadAll(bufio.NewReader(os.Stdin))

	if err != nil {
		log.Fatalln(err)
	}

	ldr.createHeap(stdin, types.Stdin, buf)

	if ldr.opts.Verbose > 1 {
		log.Println("loaded heap from stdin")
	}
}

func (ldr *Loader) createHeap(n string, t types.Heap, b []byte) {
	ldr.Lock()
	defer ldr.Unlock()

	if slices.Contains(ldr.paths, n) {
		return // already loaded
	}

	ldr.size += int64(len(b))
	ldr.paths = append(ldr.paths, n)
	ldr.heaps <- heap.New(&heap.Context{
		Name:   n,
		Type:   t,
		Limit:  ldr.opts.Limit,
		Filter: ldr.opts.Filter,
	}, b)
}

func isFilePiped(f *os.File) bool {
	fi, err := f.Stat()

	if err != nil {
		log.Fatalln(err)
	}

	return (fi.Mode() & os.ModeCharDevice) != os.ModeCharDevice
}
