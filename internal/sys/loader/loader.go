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
	"sync/atomic"

	"github.com/bmatcuk/doublestar/v4"
	"go.foxforensics.eu/fox/v4/internal/pkg"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/fox/v4/internal/sys/heap"
	"go.foxforensics.eu/fox/v4/internal/sys/memory"
)

const Stdin = "-"
const Buffer = 256
const MaxFiles = 8192
const MaxFolders = 64
const MaxArchives = 4
const MaxDeflates = 4

const (
	Initial = iota
	Extract
	Deflate
	Convert
	Format
)

type Entry struct {
	Path   string
	Part   string
	Stage  byte
	Data   []byte
	Mapped bool
}

type Options struct {
	Query    pkg.Query
	Guarded  bool
	Password string
}

type Loader struct {
	opts  *Options
	size  atomic.Uint64
	files atomic.Uint64
	paths sync.Map
	heaps chan *heap.Heap
}

func New(opts *Options) *Loader {
	return &Loader{
		opts:  opts,
		heaps: make(chan *heap.Heap, Buffer),
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

				if err = ldr.processData(ctx, &Entry{
					Path: Stdin,
					Data: bytes.TrimSpace(buf),
				}, 0); err != nil {
					slog.Error(err.Error())
				}

				break // use only stdin
			}

			select {
			case <-ctx.Done():
				return
			default:
				path, part := sys.SplitPart(path)
				ldr.loadPath(ctx, &Entry{
					Path: path,
					Part: part,
				})
			}
		}
	}()

	return ldr.heaps
}

func (ldr *Loader) Exit() {
	slog.Info(fmt.Sprintf("total size %s", sys.Humanize(ldr.size.Load())))
}

func (ldr *Loader) loadPath(ctx context.Context, e *Entry) {
	slog.Debug(fmt.Sprintf("looking for %s", e.Path))

	match, err := doublestar.FilepathGlob(e.Path)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	if len(match) == 0 {
		slog.Error(fmt.Sprintf("no files found for %s", e.Path))
		return
	}

	for _, path := range match {
		select {
		case <-ctx.Done():
			return
		default:
			fi, err := os.Stat(path)

			if err != nil {
				slog.Error(err.Error())
				continue
			}

			v := &Entry{
				Path: path,
				Part: e.Part,
			}

			if fi.IsDir() {
				err = ldr.loadDir(ctx, v, 1)
			} else {
				err = ldr.loadFile(ctx, v)
			}

			if err != nil {
				slog.Error(err.Error())
			}
		}
	}
}

func (ldr *Loader) loadDir(ctx context.Context, e *Entry, i int) error {
	if ldr.opts.Guarded && i >= MaxFolders {
		return errors.New("max folders reached")
	}

	dir, err := os.ReadDir(e.Path)

	if err != nil {
		return err
	}

	for _, f := range dir {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			v := &Entry{
				Path: filepath.Join(e.Path, f.Name()),
				Part: e.Part,
			}

			if f.IsDir() {
				err = ldr.loadDir(ctx, v, i+1)
			} else {
				err = ldr.loadFile(ctx, v)
			}

			if err != nil {
				slog.Error(err.Error())
			}
		}
	}

	return nil
}

func (ldr *Loader) loadFile(ctx context.Context, e *Entry) error {
	f, err := os.Open(e.Path)

	if err != nil {
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()

	fi, err := f.Stat()

	if err != nil {
		return err
	}

	// empty files will cause issues
	if fi.Size() == 0 {
		return ldr.createHeap(ctx, e)
	}

	e.Data, err = memory.Alloc(e.Path, f)

	if err != nil {
		return err
	}

	e.Mapped = true

	slog.Debug(fmt.Sprintf("memory mapped %s", e.Path))

	return ldr.processData(ctx, e, 0)
}

func (ldr *Loader) processData(ctx context.Context, e *Entry, i int) error {
	// check depth to protect against zip bombs
	if ldr.opts.Guarded && i >= MaxArchives {
		return errors.New("max archives reached")
	}

	// stage 1: deflate data (nested)
	for j := 1; ; j++ {
		if ldr.opts.Guarded && j >= MaxDeflates {
			return errors.New("max deflates reached")
		}

		if !ldr.deflateData(e) {
			break // no more nested
		}
	}

	// stage 2: extract data (recursive)
	if ldr.extractData(ctx, e, i) {
		return nil // ends here
	}

	// stage 3: convert data
	ldr.convertData(e)

	// stage 4: format data
	ldr.formatData(e)

	// filter for specific streams
	if len(e.Part) == 0 || strings.Contains(e.Path, e.Part) {
		return ldr.createHeap(ctx, e)
	}

	return nil
}

func (ldr *Loader) extractData(ctx context.Context, e *Entry, i int) bool {
	for _, r := range registry.Extracts {
		if r.Detect(e.Data) {
			slog.Debug(fmt.Sprintf("archive detected as possibly %s", r.Name))

			for _, s := range r.Extract(e.Data, e.Path, ldr.opts.Password) {
				slog.Debug(fmt.Sprintf("stream found %s", s.Path))

				select {
				case <-ctx.Done():
					return true
				default:
					if err := ldr.processData(ctx, &Entry{
						Path:   s.Path,
						Part:   e.Part,
						Data:   s.Data,
						Mapped: false,
						Stage:  Extract,
					}, i+1); err != nil {
						slog.Error(err.Error())
					}
				}
			}

			return true
		}
	}

	return false
}

func (ldr *Loader) deflateData(e *Entry) bool {
	for _, r := range registry.Deflates {
		if r.Detect(e.Data) {
			slog.Debug(fmt.Sprintf("deflate detected as possibly %s", r.Name))

			v, err := r.Deflate(e.Data)
			if err != nil {
				slog.Error(err.Error())
				return false
			}

			e.Data = v
			e.Mapped = false
			e.Stage = Deflate
			return true
		}
	}

	return false
}

func (ldr *Loader) convertData(e *Entry) {
	for _, r := range registry.Converts {
		if r.Detect(e.Data) {
			slog.Debug(fmt.Sprintf("convert detected as possibly %s", r.Name))

			v, err := r.Convert(e.Data)
			if err != nil {
				slog.Error(err.Error())
				return
			}

			e.Data = v
			e.Mapped = false
			e.Stage = Convert
			return
		}
	}
}

func (ldr *Loader) formatData(e *Entry) {
	for _, r := range registry.Formats {
		if r.Detect(e.Data) {
			slog.Debug(fmt.Sprintf("format detected as possibly %s", r.Name))

			v, err := r.Format(e.Data)
			if err != nil {
				slog.Error(err.Error())
				return
			}

			e.Data = v
			e.Mapped = false
			e.Stage = Format
			return
		}
	}
}

func (ldr *Loader) createHeap(ctx context.Context, e *Entry) error {
	var ok bool

	if ldr.opts.Guarded && ldr.files.Load() >= MaxFiles {
		return errors.New("max files reached")
	}

	if _, ok = ldr.paths.LoadOrStore(e.Path, pkg.Nil{}); ok {
		return nil // already loaded
	}

	ldr.size.Add(uint64(len(e.Data))) // add original size
	ldr.files.Add(1)

	if e.Data, ok = ldr.opts.Query.Reduce(e.Data); ok {
		e.Mapped = false // no more mapped
	}

	select {
	case <-ctx.Done():
		return ctx.Err()

	case ldr.heaps <- heap.New(e.Path, e.Stage, e.Mapped, e.Data):
		slog.Debug(fmt.Sprintf("loaded heap %s", e.Path))
	}

	return nil
}
