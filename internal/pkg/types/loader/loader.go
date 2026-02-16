package loader

import (
	"bufio"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/pbnjay/memory"
	"github.com/sourcegraph/conc/pool"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
	"github.com/cuhsat/fox/v4/internal/pkg/data/reader/ewf"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heap"
	"github.com/cuhsat/fox/v4/internal/pkg/types/mmap"
	"github.com/cuhsat/fox/v4/internal/pkg/types/register"
	"github.com/cuhsat/fox/v4/internal/pkg/types/share"
)

const stdin = "-"
const limit = 0.95
const peek = 8

type Options struct {
	Limit    *types.Limits
	Filter   *types.Filters
	Paths    string
	Password string
	Parallel int
	Verbose  int
	Strict   bool
	Warnings bool
}

type Loader struct {
	sync.RWMutex
	size  int64
	opts  *Options
	later []disk
	paths []string
	share *share.Share
	heaps chan *heap.Heap
}

type disk struct {
	file *os.File
	size int64
	path string
}

func New(opts *Options) *Loader {
	return &Loader{
		opts:  opts,
		heaps: make(chan *heap.Heap, opts.Parallel),
	}
}

func (ldr *Loader) Load(paths []string) <-chan *heap.Heap {
	if len(ldr.opts.Paths) > 0 {
		b, err := os.ReadFile(ldr.opts.Paths)

		if err != nil {
			log.Fatalln(err)
		}

		paths = strings.Split(strings.TrimSpace(string(b)), "\n")

		if ldr.opts.Verbose > 0 {
			for _, l := range paths {
				log.Printf("add path %s \n", l)
			}
		}
	}

	go func() {
		defer close(ldr.heaps)

		for _, path := range paths {
			if path == stdin {
				if !isPiped(os.Stdin) {
					log.Fatalln("stdin not open")
				}

				buf, err := io.ReadAll(bufio.NewReader(os.Stdin))

				if err != nil {
					log.Fatalln(err)
				}

				ldr.createData("STDIN", "", uint64(len(buf)), buf)
				break
			}

			path, part := data.SplitPart(path)

			if isRemote(path) {
				ldr.share = share.New(path)

				if ldr.opts.Verbose > 0 {
					log.Printf("use network share %s\n", ldr.share.String())
				}

				ldr.share.Mount()
			}

			_, err := os.Stat(path)

			if ldr.opts.Verbose > 0 && errors.Is(err, os.ErrNotExist) {
				log.Printf("looked for %s\n", path)
			}

			ldr.loadPath(path, part)
		}

		// combine file for ewf disks
		if path, is := isCoherent(ldr.later); is {
			ldr.combineFile(path)
		}

		// read remaining ewf disks
		for _, e := range ldr.later {
			r, err := ewf.Reader(e.file)

			if err != nil {
				log.Println(err)
			}

			ldr.createFile(e.path, "", uint64(e.size), r, e.file)
		}

		if ldr.opts.Warnings && float32(memory.FreeMemory()/memory.TotalMemory()) > limit {
			log.Println("warning: low memory may cause swapping!")
		}
	}()

	return ldr.heaps
}

func (ldr *Loader) Exit() {
	ldr.RLock()

	if ldr.share != nil {
		ldr.share.Umount()
	}

	if ldr.opts.Verbose > 0 {
		log.Printf("size %s\n", text.Humanize(ldr.size))
	}

	ldr.RUnlock()
}

func (ldr *Loader) loadPath(path, part string) {
	var root fs.FS

	base, mask := doublestar.SplitPattern(path)

	if ldr.share != nil {
		root = ldr.share.DirFS(base)
	} else {
		root = os.DirFS(base)
	}

	match, err := doublestar.Glob(root, mask)

	if err != nil {
		log.Println(err)
		return
	}

	if len(match) == 0 {
		log.Printf("file not found %s\n", path)
		return
	}

	p := pool.New().WithMaxGoroutines(ldr.opts.Parallel)

	for _, path := range match {
		path = filepath.Join(base, path)

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

	p := pool.New().WithMaxGoroutines(ldr.opts.Parallel)

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
	var f *os.File
	var err error

	if ldr.share != nil {
		f, err = ldr.share.DirFS(".").Open(path)
	} else {
		f, err = os.OpenFile(path, os.O_RDONLY, 0x400)
	}

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
		ldr.createData(path, "", 0, []byte{})
		return
	}

	// try to load the file first
	if ldr.processFile(path) {
		return
	}

	m := mmap.Map(f)

	if ldr.opts.Verbose > 2 {
		log.Printf("mapped file %s\n", path)
	}

	ldr.processData(path, part, m)
}

func (ldr *Loader) peekFile(path string) []byte {
	b := make([]byte, peek)

	f, err := os.Open(path)

	if err != nil {
		log.Fatalln(err)
	}

	r := io.LimitReader(f, peek)

	_, err = r.Read(b)

	if err != nil {
		log.Fatalln(err)
	}

	err = f.Close()

	if err != nil {
		log.Fatalln(err)
	}

	return b
}

func (ldr *Loader) combineFile(path string) {
	var file []*os.File
	var size int64

	for _, h := range ldr.later {
		file = append(file, h.file)
		size += h.size
	}

	r, err := ewf.Combine(file...)

	if err != nil {
		log.Println(err)
	}

	if ldr.opts.Verbose > 1 {
		log.Printf("disk combined as %s (%d)\n", path, len(ldr.later))
	}

	ldr.later = ldr.later[:0]

	ldr.createFile(path, "", uint64(size), r, file...)
}

func (ldr *Loader) processFile(path string) bool {
	b := ldr.peekFile(path) // peek at file header

	for _, e := range register.Readers {
		if e.Detect(b) {
			if ldr.opts.Verbose > 1 {
				log.Printf("disk detected as possibly %s\n", e.Name)
			}

			if ldr.opts.Warnings && e.Name == "ewf" {
				log.Println("warning: ewf support is experimental!")
			}

			f, err := os.Open(path)

			if err != nil {
				log.Println(err)
				break
			}

			s, err := f.Stat()

			if err != nil {
				log.Println(err)
				break
			}

			// combine ewf disks later
			if e.Name == "ewf" {
				ldr.later = append(ldr.later, disk{f, s.Size(), path})
				return true
			}

			r, err := e.Reader(f)

			if err != nil {
				log.Println(err)
				continue
			}

			ldr.createFile(path, "", uint64(s.Size()), r, f)

			return true
		}
	}

	return false
}

func (ldr *Loader) processData(path, part string, b []byte) {
	var hint string

	// 1. deflate data
	b, _ = ldr.deflateData(b)

	// 2. extract data
	if ldr.extractData(path, part, b) {
		return
	}

	// 3. convert data
	b, ok := ldr.convertData(b)

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
		ldr.createData(path, hint, uint64(len(b)), b)
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
				log.Printf("archive detected as possibly %s\n", a.Name)
			}

			p := pool.New().WithMaxGoroutines(ldr.opts.Parallel)

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

func (ldr *Loader) createData(path, hint string, size uint64, b []byte) {
	ldr.createHeap(path, size, heap.FromData(path, hint, size, b, ldr.opts.Limit))
}

func (ldr *Loader) createFile(path, hint string, size uint64, r io.ReaderAt, f ...*os.File) {
	ldr.createHeap(path, size, heap.FromFile(path, hint, size, r, f...))
}

func (ldr *Loader) createHeap(path string, size uint64, h *heap.Heap) {
	ldr.Lock()
	defer ldr.Unlock()

	if slices.Contains(ldr.paths, path) {
		return // already loaded
	}

	ldr.size += int64(size)
	ldr.paths = append(ldr.paths, path)
	ldr.heaps <- h

	if ldr.opts.Verbose > 1 {
		log.Printf("loaded heap %s\n", path)
	}
}

func isCoherent(disk []disk) (string, bool) {
	var last, name, ext string

	if len(disk) < 2 {
		return "", false
	}

	for _, d := range disk {
		ext = filepath.Ext(d.path)
		name = filepath.Base(d.path)
		name = name[0 : len(name)-len(ext)]

		if last != "" && last != name {
			return "", false
		}

		last = name
	}

	return name, true
}

func isRemote(path string) bool {
	return strings.HasPrefix("//", strings.TrimPrefix("smb:", path))
}

func isPiped(f *os.File) bool {
	fi, err := f.Stat()

	if err != nil {
		log.Fatalln(err)
	}

	return (fi.Mode() & os.ModeCharDevice) != os.ModeCharDevice
}
