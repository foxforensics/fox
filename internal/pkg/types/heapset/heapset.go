package heapset

import (
	"bufio"
	"errors"
	"fmt"
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

	szip "github.com/cuhsat/fox/v4/internal/pkg/data/archive/7z"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/ar"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/cab"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/cpio"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/rar"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/rpm"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/tar"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/xar"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/zip"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/br"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/bzip2"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/gzip"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/kanzi"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/lz4"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/lzip"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/lzw"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/minlz"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/s2"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/snappy"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/xz"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/zlib"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/zstd"
	"github.com/cuhsat/fox/v4/internal/pkg/data/parser/linux/journal"
	"github.com/cuhsat/fox/v4/internal/pkg/data/parser/windows/evtx"
	"github.com/cuhsat/fox/v4/internal/pkg/data/parser/windows/pe"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heap"
)

const stdin = "-"

type Options struct {
	Limit     *types.Limits
	Filter    *types.Filters
	Input     string
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

func (hs *HeapSet) ThrowAway() {
	var n int64

	hs.Lock()

	for _, h := range hs.heaps {
		n += h.Size()
		h.ThrowAway()
	}

	hs.heaps = hs.heaps[:0]

	hs.Unlock()

	if hs.opts.Verbose > 0 {
		log.Printf("size %s\n", human(n))
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

	if !hs.opts.NoDeflate {
		if b, ok = hs.deflate(b); ok {
			data = true
		}

		if hs.extract(path, part, b) {
			return
		}
	}

	if !hs.opts.NoConvert {
		if b, ok = hs.convert(b); ok {
			data = true
		}
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

	var fn data.Extract

	switch {
	case ar.Detect(b):
		fn = ar.Extract
	case cab.Detect(b):
		fn = cab.Extract
	case cpio.Detect(b):
		fn = cpio.Extract
	case rar.Detect(b):
		fn = rar.Extract
	case rpm.Detect(b):
		fn = rpm.Extract
	case szip.Detect(b):
		fn = szip.Extract
	case tar.Detect(b):
		fn = tar.Extract
	case xar.Detect(b):
		fn = xar.Extract
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

func (hs *HeapSet) deflate(b []byte) ([]byte, bool) {
	var fn data.Deflate

	switch {
	case br.Detect(b):
		fn = br.Deflate
	case bzip2.Detect(b):
		fn = bzip2.Deflate
	case gzip.Detect(b):
		fn = gzip.Deflate
	case kanzi.Detect(b):
		fn = kanzi.Deflate
	case lz4.Detect(b):
		fn = lz4.Deflate
	case lzip.Detect(b):
		fn = lzip.Deflate
	case lzw.Detect(b):
		fn = lzw.Deflate
	case minlz.Detect(b):
		fn = minlz.Deflate
	case s2.Detect(b):
		fn = s2.Deflate
	case snappy.Detect(b):
		fn = snappy.Deflate
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
	var fn data.Convert

	switch {
	case evtx.Detect(b):
		fn = evtx.Convert
	case journal.Detect(b):
		fn = journal.Convert
	case pe.Detect(b):
		fn = pe.Convert
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

func human(i int64) string {
	const m = int64(1024) // IEC prefix

	if i < m {
		return fmt.Sprintf("%db", i)
	}

	d, e := m, 0

	for n := i / m; n >= m; n /= m {
		d *= m
		e++
	}

	return fmt.Sprintf("%.1f%c", float64(i)/float64(d), "kmgtpezyrq"[e])
}

func debug(v any) string {
	s := runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name()
	t := strings.SplitAfter(s, "/")

	return strings.Split(t[len(t)-1], ".")[0]
}
