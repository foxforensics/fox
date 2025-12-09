package heapset

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/edsrzf/mmap-go"

	szip "github.com/cuhsat/fox/v4/internal/pkg/files/archive/7z"

	"github.com/cuhsat/fox/v4/internal/pkg/files"
	"github.com/cuhsat/fox/v4/internal/pkg/files/archive/cab"
	"github.com/cuhsat/fox/v4/internal/pkg/files/archive/rar"
	"github.com/cuhsat/fox/v4/internal/pkg/files/archive/tar"
	"github.com/cuhsat/fox/v4/internal/pkg/files/archive/zip"
	"github.com/cuhsat/fox/v4/internal/pkg/files/deflate/br"
	"github.com/cuhsat/fox/v4/internal/pkg/files/deflate/bzip2"
	"github.com/cuhsat/fox/v4/internal/pkg/files/deflate/gzip"
	"github.com/cuhsat/fox/v4/internal/pkg/files/deflate/knz"
	"github.com/cuhsat/fox/v4/internal/pkg/files/deflate/lz4"
	"github.com/cuhsat/fox/v4/internal/pkg/files/deflate/lzip"
	"github.com/cuhsat/fox/v4/internal/pkg/files/deflate/lzw"
	"github.com/cuhsat/fox/v4/internal/pkg/files/deflate/mz"
	"github.com/cuhsat/fox/v4/internal/pkg/files/deflate/s2"
	"github.com/cuhsat/fox/v4/internal/pkg/files/deflate/sz"
	"github.com/cuhsat/fox/v4/internal/pkg/files/deflate/xz"
	"github.com/cuhsat/fox/v4/internal/pkg/files/deflate/zlib"
	"github.com/cuhsat/fox/v4/internal/pkg/files/deflate/zstd"
	"github.com/cuhsat/fox/v4/internal/pkg/files/format/evtx"
	"github.com/cuhsat/fox/v4/internal/pkg/files/format/journal"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heap"
)

const Stdin = "-"

type Options struct {
	Limit     *types.Limits
	Filter    *types.Filters
	Password  string
	NoDeflate bool
	NoConvert bool
	Verbose   int
}

type HeapSet struct {
	sync.RWMutex
	opts  *Options     // set options
	heaps []*heap.Heap // set heaps
}

func New(paths []string, opts *Options) *HeapSet {
	hs := HeapSet{opts: opts}

	if isPiped(os.Stdin) {
		paths = append(paths, Stdin)
	}

	for _, path := range paths {
		if path == Stdin {
			hs.addPipe()
			return &hs
		}

		_, err := os.Stat(path)

		if errors.Is(err, os.ErrNotExist) {
			log.Printf("%s does not exist\n", path)
		}

		hs.loadPath(path)
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

func (hs *HeapSet) ThrowAway() {
	hs.Lock()

	for _, h := range hs.heaps {
		h.ThrowAway()
	}

	hs.heaps = hs.heaps[:0]

	hs.Unlock()
}

func (hs *HeapSet) loadPath(path string) {
	match, err := doublestar.FilepathGlob(path)

	if err != nil {
		log.Println(err)
		return
	}

	var wg sync.WaitGroup

	for _, m := range match {
		fi, err := os.Stat(m)

		if err != nil {
			log.Println(err)
			continue
		}

		wg.Add(1)

		go func() {
			if fi.IsDir() {
				hs.loadDir(m)
			} else {
				hs.loadFile(m)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func (hs *HeapSet) loadDir(path string) {
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
				hs.loadFile(filepath.Join(path, f.Name()))
				wg.Done()
			}()
		}
	}

	wg.Wait()
}

func (hs *HeapSet) loadFile(path string) {
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

	hs.process(path, b, false)
}

func (hs *HeapSet) process(path string, b []byte, data bool) {
	var ok bool

	if !hs.opts.NoDeflate {
		if b, ok = hs.deflate(b); ok {
			data = true
		}

		if hs.extract(path, b) {
			return
		}
	}

	if !hs.opts.NoConvert {
		if b, ok = hs.convert(b); ok {
			data = true
		}
	}

	if data {
		hs.addData(path, b)
	} else {
		hs.addFile(path, b)
	}
}

func (hs *HeapSet) extract(path string, b []byte) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Println("archive corrupt or password wrong")
			return
		}
	}()

	var fn files.Extract

	switch {
	case cab.Detect(b):
		fn = cab.Extract
	case rar.Detect(b):
		fn = rar.Extract
	case szip.Detect(b):
		fn = szip.Extract
	case tar.Detect(b):
		fn = tar.Extract
	case zip.Detect(b):
		fn = zip.Extract
	default:
		return false
	}

	if hs.opts.Verbose > 1 {
		log.Printf("format detected %s\n", debug(fn))
	}

	var wg sync.WaitGroup

	for _, e := range fn(b, path, hs.opts.Password) {
		wg.Add(1)

		go func() {
			hs.process(e.Path, e.Data, true)
			wg.Done()
		}()
	}

	wg.Wait()

	return true
}

func (hs *HeapSet) deflate(b []byte) ([]byte, bool) {
	var fn files.Deflate

	switch {
	case br.Detect(b):
		fn = br.Deflate
	case bzip2.Detect(b):
		fn = bzip2.Deflate
	case gzip.Detect(b):
		fn = gzip.Deflate
	case knz.Detect(b):
		fn = knz.Deflate
	case lz4.Detect(b):
		fn = lz4.Deflate
	case lzip.Detect(b):
		fn = lzip.Deflate
	case lzw.Detect(b):
		fn = lzw.Deflate
	case mz.Detect(b):
		fn = mz.Deflate
	case s2.Detect(b):
		fn = s2.Deflate
	case sz.Detect(b):
		fn = sz.Deflate
	case xz.Detect(b):
		fn = xz.Deflate
	case zlib.Detect(b):
		fn = zlib.Deflate
	case zstd.Detect(b):
		fn = zstd.Deflate
	default:
		return b, false
	}

	if hs.opts.Verbose > 1 {
		log.Printf("format detected %s\n", debug(fn))
	}

	r, err := fn(b)

	if err != nil {
		log.Println(err)
	}

	return r, true
}

func (hs *HeapSet) convert(b []byte) ([]byte, bool) {
	var fn files.Convert

	switch {
	case evtx.Detect(b):
		fn = evtx.Convert
	case journal.Detect(b):
		fn = journal.Convert
	default:
		return b, false
	}

	if hs.opts.Verbose > 1 {
		log.Printf("format detected %s\n", debug(fn))
	}

	r, err := fn(b)

	if err != nil {
		log.Println(err)
	}

	return r, true
}

func (hs *HeapSet) addFile(path string, b []byte) {
	hs.addHeap(path, types.Regular, b)

	if hs.opts.Verbose > 0 {
		log.Printf("loaded heap from file %s\n", path)
	}
}

func (hs *HeapSet) addData(name string, b []byte) {
	hs.addHeap(name, types.Deflate, b)

	if hs.opts.Verbose > 0 {
		log.Printf("loaded heap from data %s\n", name)
	}
}

func (hs *HeapSet) addPipe() {
	if !isPiped(os.Stdin) {
		log.Fatal("stdin not open")
	}

	buf, err := io.ReadAll(bufio.NewReader(os.Stdin))

	if err != nil {
		log.Fatal(err)
	}

	hs.addHeap(Stdin, types.Deflate, buf)

	if hs.opts.Verbose > 0 {
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
		log.Fatal(err)
	}

	return (fi.Mode() & os.ModeCharDevice) != os.ModeCharDevice
}

func debug(v any) string {
	s := runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name()
	t := strings.SplitAfter(s, "/")

	return strings.Split(t[len(t)-1], ".")[0]
}
