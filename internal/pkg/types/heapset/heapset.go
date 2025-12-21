package heapset

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/edsrzf/mmap-go"

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
	Input    string
	Password string
	Verbose  int
}

type HeapSet struct {
	sync.RWMutex
	opts  *Options     // set options
	heaps []*heap.Heap // set heaps
}

func New(paths []string, opts *Options) *HeapSet {
	hs := HeapSet{opts: opts}

	if isPiped(os.Stdin) {
		paths = append(paths, stdin)
	}

	if len(opts.Input) > 0 {
		hs.addInput([]byte(opts.Input))
	}

	for _, path := range paths {
		if path == stdin {
			hs.addStdin()
			return &hs
		}

		path, part := data.SplitPart(path)

		_, err := os.Stat(path)

		if errors.Is(err, os.ErrNotExist) {
			log.Printf("looked for %s\n", path)
		}

		hs.loadPath(path, part)
	}

	return &hs
}

func (hs *HeapSet) Len() int {
	hs.RLock()
	defer hs.RUnlock()
	return len(hs.heaps)
}

func (hs *HeapSet) Get() []*heap.Heap {
	hs.RLock()
	defer hs.RUnlock()

	r := hs.heaps[:]

	sort.SliceStable(r, func(i, j int) bool {
		return r[i].Name < r[j].Name
	})

	return r
}

func (hs *HeapSet) Discard() {
	var n int64

	hs.Lock()

	for _, h := range hs.heaps {
		n += h.Size()
		h.Discard()
	}

	hs.heaps = hs.heaps[:0]

	hs.Unlock()

	if hs.opts.Verbose > 0 {
		log.Printf("size %s\n", text.Humanize(n))
	}
}

func (hs *HeapSet) loadPath(path, part string) {
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
				hs.loadDir(path, part)
			} else {
				hs.loadFile(path, part)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func (hs *HeapSet) loadDir(path, part string) {
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
				hs.loadFile(filepath.Join(path, f.Name()), part)
				wg.Done()
			}()
		}
	}

	wg.Wait()
}

func (hs *HeapSet) loadFile(path, part string) {
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
		hs.addData(path, []byte{})
		return
	}

	b, err := mmap.Map(f, mmap.RDONLY, 0)

	if err != nil {
		log.Println(err)
		return
	}

	hs.process(path, part, b, false)
}

func (hs *HeapSet) process(path, part string, b []byte, data bool) {
	var ok bool

	// 1. deflate data
	if b, ok = hs.deflate(b); ok {
		data = true
	}

	// 2. extract data
	if hs.extract(path, part, b) {
		return
	}

	// 3. format data
	if b, ok = hs.format(b); ok {
		data = true
	}

	// filter for specific streams
	if strings.Contains(path, part) {
		if data {
			hs.addData(path, b)
		} else {
			hs.addFile(path, b)
		}
	}
}

func (hs *HeapSet) extract(path, part string, b []byte) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Println("archive corrupt or password wrong")
			return
		}
	}()

	for _, a := range register.Archives {
		if a.Detect(b) {
			if hs.opts.Verbose > 1 {
				log.Printf("archive detected %s\n", a.Name)
			}

			var wg sync.WaitGroup

			for _, e := range a.Extract(b, path, hs.opts.Password) {
				wg.Add(1)

				go func() {
					if hs.opts.Verbose > 2 {
						log.Printf("stream detected %s\n", e.Path)
					}

					hs.process(e.Path, part, e.Data, true)

					wg.Done()
				}()
			}

			wg.Wait()

			return true
		}
	}

	return false
}

func (hs *HeapSet) deflate(b []byte) ([]byte, bool) {
	for _, d := range register.Deflates {
		if d.Detect(b) {
			if hs.opts.Verbose > 1 {
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

func (hs *HeapSet) format(b []byte) ([]byte, bool) {
	for _, f := range register.Formats {
		if f.Detect(b) {
			if hs.opts.Verbose > 1 {
				log.Printf("format detected %s\n", f.Name)
			}

			r, err := f.Format(b, hs.opts.Verbose)

			if err != nil {
				log.Println(err)
			}

			return r, true
		}
	}

	return b, false
}

func (hs *HeapSet) addFile(path string, b []byte) {
	hs.addHeap(path, types.Regular, b)

	if hs.opts.Verbose > 1 {
		log.Printf("loaded heap from file %s\n", path)
	}
}

func (hs *HeapSet) addData(name string, b []byte) {
	hs.addHeap(name, types.Deflate, b)

	if hs.opts.Verbose > 1 {
		log.Printf("loaded heap from data %s\n", name)
	}
}

func (hs *HeapSet) addInput(b []byte) {
	hs.addHeap("input", types.String, b)

	if hs.opts.Verbose > 1 {
		log.Println("loaded heap from input")
	}
}

func (hs *HeapSet) addStdin() {
	if !isPiped(os.Stdin) {
		log.Fatalln("stdin not open")
	}

	buf, err := io.ReadAll(bufio.NewReader(os.Stdin))

	if err != nil {
		log.Fatalln(err)
	}

	hs.addHeap(stdin, types.Stdin, buf)

	if hs.opts.Verbose > 1 {
		log.Println("loaded heap from stdin")
	}
}

func (hs *HeapSet) addHeap(s string, t types.Heap, b []byte) {
	hs.Lock()
	defer hs.Unlock()

	for _, h := range hs.heaps {
		if h.Name == s {
			return // already loaded
		}
	}

	hs.heaps = append(hs.heaps, heap.New(
		&heap.Context{
			Name:   s,
			Type:   t,
			Limit:  hs.opts.Limit,
			Filter: hs.opts.Filter,
		}, b,
	))
}

func isPiped(f *os.File) bool {
	fi, err := f.Stat()

	if err != nil {
		log.Fatalln(err)
	}

	return (fi.Mode() & os.ModeCharDevice) != os.ModeCharDevice
}
