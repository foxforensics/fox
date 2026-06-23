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
	"sync/atomic"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/sourcegraph/conc/pool"
	"go.foxforensics.eu/fox/v4/internal/pkg/types"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/fox/v4/internal/sys/heap"
	"go.foxforensics.eu/fox/v4/internal/sys/mmap"
)

const Stdin = "-"
const MaxFiles = 8192
const MaxFolders = 64
const MaxArchives = 4

type Options struct {
	Limits   *types.Limits
	Filters  *types.Filters
	Password string
	Threads  int
	Strict   bool
}

type Loader struct {
	opts  *Options
	size  atomic.Uint64
	paths *types.Set
	heaps chan *heap.Heap
}

func New(opts *Options) *Loader {
	return &Loader{
		opts:  opts,
		paths: types.NewSet(),
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
	slog.Debug(fmt.Sprintf("total size %s", sys.Humanize(ldr.size.Load())))
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
				return ldr.loadDir(ctx, path, part, 1)
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

func (ldr *Loader) loadDir(ctx context.Context, path, part string, i int) error {
	if ldr.opts.Strict && i > MaxFolders {
		return errors.New("max folders reached")
	}

	dir, err := os.ReadDir(path)

	if err != nil {
		return err
	}

	p := pool.New().
		WithContext(ctx).
		WithFirstError().
		WithMaxGoroutines(ldr.opts.Threads)

	for _, f := range dir {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if f.IsDir() {
				p.Go(func(ctx context.Context) error {
					return ldr.loadDir(ctx, filepath.Join(path, f.Name()), part, i+1)
				})
			} else {
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
	if ldr.opts.Strict && i > MaxArchives {
		return errors.New("max archives reached")
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
	if ldr.paths.Has(path) {
		return nil // already loaded
	}

	// check files to protect against zip bombs
	if ldr.opts.Strict && ldr.paths.Len() >= MaxFiles {
		return errors.New("max files reached")
	}

	// add original size
	ldr.size.Add(uint64(len(b)))

	b = ldr.opts.Limits.Reduce(b)

	ldr.paths.Set(path)
	ldr.heaps <- heap.New(path, hint, b)

	slog.Debug(fmt.Sprintf("loaded heap %s", path))

	return nil
}
