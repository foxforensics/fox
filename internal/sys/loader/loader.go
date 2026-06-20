package loader

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/sourcegraph/conc/pool"
	_zip "go.foxforensics.eu/fox/v4/internal/pkg/archive/7z"
	"go.foxforensics.eu/fox/v4/internal/pkg/archive/ar"
	"go.foxforensics.eu/fox/v4/internal/pkg/archive/cab"
	"go.foxforensics.eu/fox/v4/internal/pkg/archive/cpio"
	"go.foxforensics.eu/fox/v4/internal/pkg/archive/iso"
	"go.foxforensics.eu/fox/v4/internal/pkg/archive/msi"
	"go.foxforensics.eu/fox/v4/internal/pkg/archive/rar"
	"go.foxforensics.eu/fox/v4/internal/pkg/archive/rpm"
	"go.foxforensics.eu/fox/v4/internal/pkg/archive/tar"
	"go.foxforensics.eu/fox/v4/internal/pkg/archive/xar"
	"go.foxforensics.eu/fox/v4/internal/pkg/archive/zip"
	"go.foxforensics.eu/fox/v4/internal/pkg/binary/bin/elf"
	"go.foxforensics.eu/fox/v4/internal/pkg/binary/bin/ese"
	"go.foxforensics.eu/fox/v4/internal/pkg/binary/bin/lnk"
	"go.foxforensics.eu/fox/v4/internal/pkg/binary/bin/mft"
	"go.foxforensics.eu/fox/v4/internal/pkg/binary/bin/pe"
	"go.foxforensics.eu/fox/v4/internal/pkg/binary/bin/pf"
	"go.foxforensics.eu/fox/v4/internal/pkg/binary/bin/pst"
	"go.foxforensics.eu/fox/v4/internal/pkg/binary/log/evtx"
	"go.foxforensics.eu/fox/v4/internal/pkg/binary/log/fortinet"
	"go.foxforensics.eu/fox/v4/internal/pkg/binary/log/journal"
	"go.foxforensics.eu/fox/v4/internal/pkg/deflate/bgzf"
	"go.foxforensics.eu/fox/v4/internal/pkg/deflate/br"
	"go.foxforensics.eu/fox/v4/internal/pkg/deflate/bzip2"
	"go.foxforensics.eu/fox/v4/internal/pkg/deflate/gzip"
	"go.foxforensics.eu/fox/v4/internal/pkg/deflate/kanzi"
	"go.foxforensics.eu/fox/v4/internal/pkg/deflate/lz4"
	"go.foxforensics.eu/fox/v4/internal/pkg/deflate/lzfse"
	"go.foxforensics.eu/fox/v4/internal/pkg/deflate/lzip"
	"go.foxforensics.eu/fox/v4/internal/pkg/deflate/lznt1"
	"go.foxforensics.eu/fox/v4/internal/pkg/deflate/lzo"
	"go.foxforensics.eu/fox/v4/internal/pkg/deflate/lzw"
	"go.foxforensics.eu/fox/v4/internal/pkg/deflate/minlz"
	"go.foxforensics.eu/fox/v4/internal/pkg/deflate/s2"
	"go.foxforensics.eu/fox/v4/internal/pkg/deflate/snappy"
	"go.foxforensics.eu/fox/v4/internal/pkg/deflate/xz"
	"go.foxforensics.eu/fox/v4/internal/pkg/deflate/zlib"
	"go.foxforensics.eu/fox/v4/internal/pkg/deflate/zstd"
	"go.foxforensics.eu/fox/v4/internal/pkg/format/json"
	"go.foxforensics.eu/fox/v4/internal/pkg/format/jsonl"
	"go.foxforensics.eu/fox/v4/internal/pkg/format/xml"
	"go.foxforensics.eu/fox/v4/internal/pkg/types"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/fox/v4/internal/sys/heap"
	"go.foxforensics.eu/fox/v4/internal/sys/mmap"
	"go.foxforensics.eu/fox/v4/internal/sys/register"
)

const Stdin = "-"
const MaxDepth = 3
const MaxFiles = 10000

var registry = register.Registry

type Options struct {
	Limits   *types.Limits
	Filters  *types.Filters
	Password string
	Threads  int
	Strict   bool
}

type Loader struct {
	sync.RWMutex
	size  uint64
	opts  *Options
	paths types.Set
	heaps chan *heap.Heap
}

func New(opts *Options) *Loader {
	return &Loader{
		opts:  opts,
		paths: make(types.Set),
		heaps: make(chan *heap.Heap, opts.Threads),
	}
}

func (ldr *Loader) Load(ctx context.Context, paths []string) <-chan *heap.Heap {
	go func() {
		defer close(ldr.heaps)

		for _, path := range paths {
			// read file content from stdin
			if path == Stdin {
				fi, err := os.Stdin.Stat()

				if err != nil {
					slog.Error(err.Error())
					continue
				}

				if (fi.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
					slog.Error("stdin is not open")
					continue
				}

				buf, err := io.ReadAll(bufio.NewReader(os.Stdin))

				if err != nil {
					slog.Error(err.Error())
					continue
				}

				err = ldr.processData(Stdin, "", bytes.TrimSpace(buf), 0)

				if err != nil {
					slog.Error(err.Error())
				}

				break // use only stdin
			}

			select {
			case <-ctx.Done():
				return
			default:
				p1, p2 := sys.SplitPart(path)
				ldr.loadPath(ctx, p1, p2)
			}
		}
	}()

	return ldr.heaps
}

func (ldr *Loader) Exit() {
	ldr.RLock()
	slog.Debug(fmt.Sprintf("total size %s", sys.Humanize(ldr.size)))
	ldr.RUnlock()
}

func (ldr *Loader) loadPath(ctx context.Context, path, part string) {
	v, err := filepath.Abs(path)

	if err == nil {
		path = v
	}

	slog.Debug(fmt.Sprintf("looking for %s", path))

	match, err := doublestar.FilepathGlob(path)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	if len(match) == 0 {
		slog.Error(fmt.Sprintf("no files found for %s", path))
		return
	}

	p := pool.New().
		WithContext(ctx).
		WithFirstError().
		WithMaxGoroutines(ldr.opts.Threads)

	for _, path := range match {
		fi, err := os.Stat(path)

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		p.Go(func(ctx context.Context) error {
			if fi.IsDir() {
				return ldr.loadDir(ctx, path, part)
			} else {
				return ldr.loadFile(ctx, path, part)
			}
		})
	}

	if err = p.Wait(); err != nil {
		if errors.Is(err, context.Canceled) {
			slog.Info("canceled")
		} else {
			slog.Error(err.Error())
		}
	}
}

func (ldr *Loader) loadDir(ctx context.Context, path, part string) error {
	dir, err := os.ReadDir(path)

	if err != nil {
		return err
	}

	p := pool.New().
		WithContext(ctx).
		WithFirstError().
		WithMaxGoroutines(ldr.opts.Threads)

	for _, f := range dir {
		if !f.IsDir() {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				p.Go(func(ctx context.Context) error {
					return ldr.loadFile(ctx, filepath.Join(path, f.Name()), part)
				})
			}
		}
	}

	return p.Wait()
}

func (ldr *Loader) loadFile(ctx context.Context, path, part string) error {
	f, err := os.Open(path)

	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	fi, err := f.Stat()

	if err != nil {
		return err
	}

	// empty files will cause issues
	if fi.Size() == 0 {
		return ldr.createHeap(path, "", []byte{})
	}

	b, err := mmap.Map(f)

	if err != nil {
		return err
	}

	slog.Debug(fmt.Sprintf("loaded file %s", path))

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return ldr.processData(path, part, b, 0)
	}
}

func (ldr *Loader) processData(path, part string, b []byte, i int) error {
	var hint string
	var ok bool

	// check depth to protect against zip bombs
	if ldr.opts.Strict && i > MaxDepth {
		return errors.New("max depth reached")
	}

	// 1. deflate data
	for {
		if b, ok = ldr.deflateData(b); !ok {
			break
		}
	}

	// 2. extract data (recursive)
	if ldr.extractData(path, part, b, i) {
		return nil
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
		return ldr.createHeap(path, hint, b)
	}

	return nil
}

func (ldr *Loader) extractData(path, part string, b []byte, i int) bool {
	defer func() {
		// most libraries can not differentiate between invalid data and wrong passwords
		if err := recover(); err != nil {
			slog.Error("archive corrupt or password wrong")
			return
		}
	}()

	for _, a := range registry.Extracts {
		if a.Detect(b) {
			slog.Debug(fmt.Sprintf("archive detected as possibly %s", a.Name))

			p := pool.New().
				WithErrors().
				WithFirstError().
				WithMaxGoroutines(ldr.opts.Threads)

			for _, e := range a.Extract(b, path, ldr.opts.Password) {
				p.Go(func() error {
					slog.Debug(fmt.Sprintf("stream detected as possibly %s", e.Path))
					return ldr.processData(e.Path, part, e.Data, i+1)
				})
			}

			if err := p.Wait(); err != nil {
				slog.Error(err.Error())
			}

			return true
		}
	}

	return false
}

func (ldr *Loader) deflateData(b []byte) ([]byte, bool) {
	for _, e := range registry.Deflates {
		if e.Detect(b) {
			slog.Debug(fmt.Sprintf("deflate detected as possibly %s", e.Name))

			r, err := e.Deflate(b)

			if err != nil {
				slog.Error(err.Error())

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
	for _, e := range registry.Converts {
		if e.Detect(b) {
			slog.Debug(fmt.Sprintf("convert detected as possibly %s", e.Name))

			r, err := e.Convert(b)

			if err != nil {
				slog.Error(err.Error())

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
	for _, e := range registry.Formats {
		if e.Detect(b) {
			slog.Debug(fmt.Sprintf("format detected as possibly %s", e.Name))

			r, err := e.Format(b)

			if err != nil {
				slog.Error(err.Error())

				if ldr.opts.Strict {
					r = b // ignore partly result
				}
			}

			return r, true
		}
	}

	return b, false
}

func (ldr *Loader) createHeap(path, hint string, b []byte) error {
	ldr.Lock()
	defer ldr.Unlock()

	if ldr.paths.Has(path) {
		return nil // already loaded
	}

	// check files to protect against zip bombs
	if ldr.opts.Strict && len(ldr.paths) >= MaxFiles {
		return errors.New("max files reached")
	}

	// add original size
	ldr.size += uint64(len(b))

	b = ldr.opts.Limits.Reduce(b)

	ldr.paths.Set(path)
	ldr.heaps <- heap.New(path, hint, b)

	slog.Debug(fmt.Sprintf("loaded heap %s", path))

	return nil
}

func RegisterDeflates() {
	register.Deflate("bgzf", bgzf.Detect, bgzf.Deflate)
	register.Deflate("br", br.Detect, br.Deflate)
	register.Deflate("bzip2", bzip2.Detect, bzip2.Deflate)
	register.Deflate("gzip", gzip.Detect, gzip.Deflate)
	register.Deflate("kanzi", kanzi.Detect, kanzi.Deflate)
	register.Deflate("lz4", lz4.Detect, lz4.Deflate)
	register.Deflate("lzip", lzip.Detect, lzip.Deflate)
	register.Deflate("lzo", lzo.Detect, lzo.Deflate)
	register.Deflate("lzfse", lzfse.Detect, lzfse.Deflate)
	register.Deflate("lznt1", lznt1.Detect, lznt1.Deflate)
	register.Deflate("lzw", lzw.Detect, lzw.Deflate)
	register.Deflate("minlz", minlz.Detect, minlz.Deflate)
	register.Deflate("s2", s2.Detect, s2.Deflate)
	register.Deflate("snappy", snappy.Detect, snappy.Deflate)
	register.Deflate("xz", xz.Detect, xz.Deflate)
	register.Deflate("zlib", zlib.Detect, zlib.Deflate)
	register.Deflate("zstd", zstd.Detect, zstd.Deflate)
}

func RegisterExtracts() {
	register.Extract("7z", _zip.Detect, _zip.Extract)
	register.Extract("ar", ar.Detect, ar.Extract)
	register.Extract("cab", cab.Detect, cab.Extract)
	register.Extract("cpio", cpio.Detect, cpio.Extract)
	register.Extract("iso", iso.Detect, iso.Extract)
	register.Extract("msi", msi.Detect, msi.Extract)
	register.Extract("rar", rar.Detect, rar.Extract)
	register.Extract("rpm", rpm.Detect, rpm.Extract)
	register.Extract("tar", tar.Detect, tar.Extract)
	register.Extract("xar", xar.Detect, xar.Extract)
	register.Extract("zip", zip.Detect, zip.Extract)
}

func RegisterConverts() {
	register.Convert("elf", elf.Detect, elf.Convert)
	register.Convert("ese", ese.Detect, ese.Convert)
	register.Convert("lnk", lnk.Detect, lnk.Convert)
	register.Convert("mft", mft.Detect, mft.Convert)
	register.Convert("pe", pe.Detect, pe.Convert)
	register.Convert("pf", pf.Detect, pf.Convert)
	register.Convert("pst", pst.Detect, pst.Convert)
	register.Convert("evtx", evtx.Detect, evtx.Convert)
	register.Convert("fortinet", fortinet.Detect, fortinet.Convert)
	register.Convert("journal", journal.Detect, journal.Convert)
}

func RegisterFormats() {
	register.Format("json", json.Detect, json.Format)
	register.Format("jsonl", jsonl.Detect, jsonl.Format)
	register.Format("xml", xml.Detect, xml.Format)
}
