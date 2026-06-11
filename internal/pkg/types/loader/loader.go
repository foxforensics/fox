package loader

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/sourcegraph/conc/pool"

	"go.foxforensics.dev/fox/v4/internal/pkg/file"
	"go.foxforensics.dev/fox/v4/internal/pkg/text"
	"go.foxforensics.dev/fox/v4/internal/pkg/types"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/heap"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/mmap"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/register"
)

const Stdin = "-"

type Options struct {
	Limits   *types.Limits
	Filters  *types.Filters
	Password string
	Threads  int
	Verbose  int
	Strict   bool
}

type Loader struct {
	sync.RWMutex
	size  int64
	opts  *Options
	paths []string
	heaps chan *heap.Heap
}

func New(opts *Options) *Loader {
	return &Loader{
		opts:  opts,
		heaps: make(chan *heap.Heap, opts.Threads),
	}
}

func (ldr *Loader) Load(paths []string) <-chan *heap.Heap {
	go func() {
		defer close(ldr.heaps)

		for _, path := range paths {
			// read file content from stdin
			if path == Stdin {
				if !isPiped(os.Stdin) {
					log.Fatalln("stdin not open")
				}

				buf, err := io.ReadAll(bufio.NewReader(os.Stdin))

				if err != nil {
					log.Fatalln(err)
				}

				ldr.processData(Stdin, "", bytes.TrimSpace(buf))
				break
			}

			ldr.loadPath(file.SplitPart(path))
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
	if ldr.opts.Verbose > 0 {
		log.Printf("looking for %s\n", path)
	}

	match, err := doublestar.FilepathGlob(filepath.Clean(path))

	if err != nil {
		log.Println(err)
		return
	}

	if len(match) == 0 {
		log.Printf("no files found for %s\n", path)
		return
	}

	p := pool.New().WithMaxGoroutines(ldr.opts.Threads)

	for _, path := range match {
		fi, err := os.Stat(path)

		if err != nil {
			log.Println(err)
			continue
		}

		p.Go(func() {
			if fi.IsDir() {
				ldr.loadDir(path, part)
			} else {
				ldr.loadFile(path, part)
			}
		})
	}

	p.Wait()
}

func (ldr *Loader) loadDir(path, part string) {
	dir, err := os.ReadDir(path)

	if err != nil {
		log.Println(err)
		return
	}

	p := pool.New().WithMaxGoroutines(ldr.opts.Threads)

	for _, f := range dir {
		if !f.IsDir() {
			p.Go(func() {
				ldr.loadFile(filepath.Join(path, f.Name()), part)
			})
		}
	}

	p.Wait()
}

func (ldr *Loader) loadFile(path, part string) {
	f, err := os.Open(path)

	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		_ = f.Close()
	}()

	fi, err := f.Stat()

	if err != nil {
		log.Println(err)
		return
	}

	// empty files will cause issues
	if fi.Size() == 0 {
		ldr.createHeap(path, "", 0, []byte{})
		return
	}

	b := mmap.Map(f)

	if ldr.opts.Verbose > 2 {
		log.Printf("mapped file %s\n", path)
	}

	ldr.processData(path, part, b)
}

func (ldr *Loader) processData(path, part string, b []byte) {
	var hint string
	var ok bool

	// 1. deflate data (recursive)
	for {
		if b, ok = ldr.deflateData(b); !ok {
			break
		}
	}

	// 2. extract data (multiply)
	if ldr.extractData(path, part, b) {
		return
	}

	// 3. convert data
	b, ok = ldr.convertData(b)

	// default conversion format
	if ok {
		hint = "json"
	}

	// 4. format data
	b, ok = ldr.formatData(b)

	// only formating style
	if ok {
		hint = "json"
	}

	// filter for specific streams
	if strings.Contains(path, part) {
		ldr.createHeap(path, hint, uint64(len(b)), b)
	}
}

func (ldr *Loader) extractData(path, part string, b []byte) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Println("archive corrupt or password wrong")
			return
		}
	}()

	for _, a := range register.Extracts {
		if a.Detect(b) {
			if ldr.opts.Verbose > 1 {
				log.Printf("archive detected as possibly %s\n", a.Name)
			}

			p := pool.New().WithMaxGoroutines(ldr.opts.Threads)

			for _, e := range a.Extract(b, path, ldr.opts.Password) {
				p.Go(func() {
					if ldr.opts.Verbose > 2 {
						log.Printf("stream detected as %s\n", e.Path)
					}

					ldr.processData(e.Path, part, e.Data)
				})
			}

			p.Wait()

			return true
		}
	}

	return false
}

func (ldr *Loader) deflateData(b []byte) ([]byte, bool) {
	for _, e := range register.Deflates {
		if e.Detect(b) {
			if ldr.opts.Verbose > 1 {
				log.Printf("deflate detected as possibly %s\n", e.Name)
			}

			r, err := e.Deflate(b)

			if err != nil {
				log.Println(err)

				if ldr.opts.Strict {
					r = b // ignore partly result
				}
			}

			return r, true
		}
	}

	return b, false
}

func (ldr *Loader) convertData(b []byte) ([]byte, bool) {
	for _, e := range register.Converts {
		if e.Detect(b) {
			if ldr.opts.Verbose > 1 {
				log.Printf("convert detected as possibly %s\n", e.Name)
			}

			r, err := e.Convert(b)

			if err != nil {
				log.Println(err)

				if ldr.opts.Strict {
					r = b // ignore partly result
				}
			}

			return r, true
		}
	}

	return b, false
}

func (ldr *Loader) formatData(b []byte) ([]byte, bool) {
	for _, e := range register.Formats {
		if e.Detect(b) {
			if ldr.opts.Verbose > 1 {
				log.Printf("format detected as possibly %s\n", e.Name)
			}

			r, err := e.Format(b)

			if err != nil {
				log.Println(err)

				if ldr.opts.Strict {
					r = b // ignore partly result
				}
			}

			return r, true
		}
	}

	return b, false
}

func (ldr *Loader) createHeap(path, hint string, size uint64, b []byte) {
	ldr.Lock()
	defer ldr.Unlock()

	if slices.Contains(ldr.paths, path) {
		return // already loaded
	}

	b = ldr.opts.Limits.Reduce(b)

	ldr.size += int64(size)
	ldr.paths = append(ldr.paths, path)
	ldr.heaps <- heap.New(path, hint, size, b)

	if ldr.opts.Verbose > 1 {
		log.Printf("loaded heap %s\n", path)
	}
}

func isPiped(f *os.File) bool {
	fi, err := f.Stat()

	if err != nil {
		log.Println(err)
		return false
	}

	return (fi.Mode() & os.ModeCharDevice) != os.ModeCharDevice
}
