package loader

import (
	"bufio"
	"io"
	"io/fs"
	"log"
	"maps"
	"os"
	"path/filepath"
	"regexp"
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

const (
	stdin = "-"
	limit = 0.95
	peek  = 16
)

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
	size   int64
	opts   *Options
	later  []disk
	paths  []string
	heaps  chan *heap.Heap
	mounts map[string]*share.Share
}

type disk struct {
	file types.File
	size int64
	path string
}

func New(opts *Options) *Loader {
	return &Loader{
		opts:   opts,
		heaps:  make(chan *heap.Heap, opts.Parallel),
		mounts: make(map[string]*share.Share),
	}
}

func (ldr *Loader) Load(paths []string) <-chan *heap.Heap {
	// read paths from given file
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
			// read file content from stdin
			if path == stdin {
				if !isPiped(os.Stdin) {
					log.Fatalln("stdin not open")
				}

				buf, err := io.ReadAll(bufio.NewReader(os.Stdin))

				if err != nil {
					log.Fatalln(err)
				}

				ldr.processData("STDIN", "", buf)
				break
			}

			// split file and stream
			path, part := data.SplitPart(path)

			// mount remote share
			if isRemote(path) {
				mnt, tmp := share.New(path)

				ldr.mounts[mnt.String()] = mnt

				if ldr.opts.Verbose > 0 {
					log.Printf("mount share %s\n", mnt.String())
				}

				mnt.Mount()

				path = tmp
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

	for _, mnt := range maps.All(ldr.mounts) {
		if ldr.opts.Verbose > 0 {
			log.Printf("umount share %s\n", mnt.String())
		}

		mnt.Umount()
	}

	if ldr.opts.Verbose > 0 {
		log.Printf("size %s\n", text.Humanize(ldr.size))
	}

	ldr.RUnlock()
}

func (ldr *Loader) loadPath(path, part string) {
	var root fs.FS

	if ldr.opts.Verbose > 0 {
		log.Printf("looking for %s\n", path)
	}

	base, mask := doublestar.SplitPattern(path)

	if shr, ok := ldr.mounts[path]; ok {
		root = shr.DirFS(".")
	} else {
		root = os.DirFS(base)
	}

	match, err := doublestar.Glob(root, mask)

	if err != nil {
		log.Println(err)
		return
	}

	if len(match) == 0 {
		log.Printf("no files found %s\n", path)
		return
	}

	p := pool.New().WithMaxGoroutines(ldr.opts.Parallel)

	for _, path := range match {
		fi, err := fs.Stat(root, path)

		if err != nil {
			log.Println(err)
			continue
		}

		p.Go(func() {
			if fi.IsDir() {
				ldr.loadDir(root, path, part)
			} else {
				ldr.loadFile(root, path, part)
			}
		})
	}

	p.Wait()
}

func (ldr *Loader) loadDir(root fs.FS, path, part string) {
	dir, err := fs.ReadDir(root, path)

	if err != nil {
		log.Println(err)
		return
	}

	p := pool.New().WithMaxGoroutines(ldr.opts.Parallel)

	for _, f := range dir {
		if !f.IsDir() {
			p.Go(func() {
				ldr.loadFile(root, filepath.Join(path, f.Name()), part)
			})
		}
	}

	p.Wait()
}

func (ldr *Loader) loadFile(root fs.FS, path, part string) {
	f, err := root.Open(path)

	if err != nil {
		log.Println(err)
		return
	}

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
	if ldr.processFile(path, f.(types.File)) {
		return
	}

	var b []byte

	switch v := f.(type) {
	// memory map local files
	case *os.File:
		b = mmap.Map(v)

		if ldr.opts.Verbose > 2 {
			log.Printf("mapped file %s\n", path)
		}

	// read remote file whole
	default:
		b, err = fs.ReadFile(root, path)

		if err != nil {
			log.Println(err)
		}

		if ldr.opts.Verbose > 2 {
			log.Printf("loaded file %s\n", path)
		}
	}

	_ = f.Close()

	ldr.processData(path, part, b)
}

func (ldr *Loader) peekFile(file types.File) []byte {
	b := make([]byte, peek)

	_, err := file.Read(b)

	if err != nil {
		log.Fatalln(err)
	}

	_, err = file.Seek(0, io.SeekStart)

	if err != nil {
		log.Fatalln(err)
	}

	return b
}

func (ldr *Loader) combineFile(path string) {
	var f []types.File
	var c []io.Closer
	var s int64

	for _, d := range ldr.later {
		f = append(f, d.file)
		c = append(c, d.file)
		s += d.size
	}

	r, err := ewf.Combine(f...)

	if err != nil {
		log.Println(err)
	}

	if ldr.opts.Verbose > 1 {
		log.Printf("disk combined as %s (%d)\n", path, len(ldr.later))
	}

	ldr.later = ldr.later[:0]

	ldr.createFile(path, "", uint64(s), r, c...)
}

func (ldr *Loader) processFile(path string, file types.File) bool {
	b := ldr.peekFile(file) // peek at file header

	for _, e := range register.Readers {
		if e.Detect(b) {
			if ldr.opts.Verbose > 1 {
				log.Printf("disk detected as possibly %s\n", e.Name)
			}

			if ldr.opts.Warnings && e.Name == "ewf" {
				log.Println("warning: ewf support is experimental!")
			}

			s, err := file.Stat()

			if err != nil {
				log.Println(err)
				break
			}

			// combine ewf disks later
			if e.Name == "ewf" {
				ldr.later = append(ldr.later, disk{file, s.Size(), path})
				return true
			}

			r, err := e.Reader(file)

			if err != nil {
				log.Println(err)
				continue
			}

			ldr.createFile(path, "", uint64(s.Size()), r, file)

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

func (ldr *Loader) createFile(path, hint string, size uint64, r io.ReaderAt, c ...io.Closer) {
	ldr.createHeap(path, size, heap.FromFile(path, hint, size, r, c...))
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
	return regexp.MustCompile(`^\S*(smb:)?(//|\\).+`).MatchString(path)
}

func isPiped(f *os.File) bool {
	fi, err := f.Stat()

	if err != nil {
		log.Fatalln(err)
	}

	return (fi.Mode() & os.ModeCharDevice) != os.ModeCharDevice
}
