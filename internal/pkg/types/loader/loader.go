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
	paths []string
	size  int64
	ch    chan *heap.Heap
}

func New(opts *Options) *Loader {
	return &Loader{
		opts: opts,
		ch:   make(chan *heap.Heap, opts.Queue),
	}
}

func (l *Loader) Load(paths []string) <-chan *heap.Heap {
	go func() {
		defer close(l.ch)

		if isFilePiped(os.Stdin) {
			paths = append(paths, stdin)
		}

		if len(l.opts.Input) > 0 {
			l.createFromInput([]byte(l.opts.Input))
		}

		for _, path := range paths {
			if path == stdin {
				l.createFromStdin()
				break
			}

			path, part := data.SplitPart(path)

			_, err := os.Stat(path)

			if l.opts.Verbose > 0 && errors.Is(err, os.ErrNotExist) {
				log.Printf("looked for %s\n", path)
			}

			l.loadPath(path, part)
		}
	}()

	return l.ch
}

func (l *Loader) Exit() {
	l.RLock()

	if l.opts.Verbose > 0 {
		log.Printf("size %s\n", text.Humanize(l.size))
	}

	l.RUnlock()
}

func (l *Loader) loadPath(path, part string) {
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
				l.loadDir(path, part)
			} else {
				l.loadFile(path, part)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func (l *Loader) loadDir(path, part string) {
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
				l.loadFile(filepath.Join(path, f.Name()), part)
				wg.Done()
			}()
		}
	}

	wg.Wait()
}

func (l *Loader) loadFile(path, part string) {
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
		l.createFromData(path, []byte{})
		return
	}

	b, err := mmap.Map(f, mmap.RDONLY, 0)

	if err != nil {
		log.Println(err)
		return
	}

	l.processData(path, part, b, false)
}

func (l *Loader) processData(path, part string, b []byte, data bool) {
	var ok bool

	// 1. deflate data
	if b, ok = l.deflateData(b); ok {
		data = true
	}

	// 2. extract data
	if l.extractData(path, part, b) {
		return
	}

	// 3. format data
	if b, ok = l.formatData(b); ok {
		data = true
	}

	// filter for specific streams
	if strings.Contains(path, part) {
		if data {
			l.createFromData(path, b)
		} else {
			l.createFromFile(path, b)
		}
	}
}

func (l *Loader) extractData(path, part string, b []byte) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Println("archive corrupt or password wrong")
			return
		}
	}()

	for _, a := range register.Archives {
		if a.Detect(b) {
			if l.opts.Verbose > 1 {
				log.Printf("archive detected %s\n", a.Name)
			}

			var wg sync.WaitGroup

			for _, e := range a.Extract(b, path, l.opts.Password) {
				wg.Add(1)

				go func() {
					if l.opts.Verbose > 2 {
						log.Printf("stream detected %s\n", e.Path)
					}

					l.processData(e.Path, part, e.Data, true)

					wg.Done()
				}()
			}

			wg.Wait()

			return true
		}
	}

	return false
}

func (l *Loader) deflateData(b []byte) ([]byte, bool) {
	for _, d := range register.Deflates {
		if d.Detect(b) {
			if l.opts.Verbose > 1 {
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

func (l *Loader) formatData(b []byte) ([]byte, bool) {
	for _, c := range register.Converts {
		if c.Detect(b) {
			if l.opts.Verbose > 1 {
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

func (l *Loader) createFromFile(path string, b []byte) {
	l.createHeap(path, types.Regular, b)

	if l.opts.Verbose > 1 {
		log.Printf("loaded heap from file %s\n", path)
	}
}

func (l *Loader) createFromData(name string, b []byte) {
	l.createHeap(name, types.Deflate, b)

	if l.opts.Verbose > 1 {
		log.Printf("loaded heap from data %s\n", name)
	}
}

func (l *Loader) createFromInput(b []byte) {
	l.createHeap("input", types.Defined, b)

	if l.opts.Verbose > 1 {
		log.Println("loaded heap from input")
	}
}

func (l *Loader) createFromStdin() {
	if !isFilePiped(os.Stdin) {
		log.Fatalln("stdin not open")
	}

	buf, err := io.ReadAll(bufio.NewReader(os.Stdin))

	if err != nil {
		log.Fatalln(err)
	}

	l.createHeap(stdin, types.Stdin, buf)

	if l.opts.Verbose > 1 {
		log.Println("loaded heap from stdin")
	}
}

func (l *Loader) createHeap(n string, t types.Heap, b []byte) {
	l.Lock()
	defer l.Unlock()

	if slices.Contains(l.paths, n) {
		return // already loaded
	}

	l.paths = append(l.paths, n)
	l.size += int64(len(b))
	l.ch <- heap.New(&heap.Context{
		Name:   n,
		Type:   t,
		Limit:  l.opts.Limit,
		Filter: l.opts.Filter,
	}, b)
}

func isFilePiped(f *os.File) bool {
	fi, err := f.Stat()

	if err != nil {
		log.Fatalln(err)
	}

	return (fi.Mode() & os.ModeCharDevice) != os.ModeCharDevice
}
