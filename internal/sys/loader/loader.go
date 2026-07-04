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

type Options struct {
	Query    pkg.Query
	Protect  bool
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

				h := heap.FromData(Stdin, bytes.TrimSpace(buf))

				if err = ldr.processData(ctx, h, 0); err != nil {
					slog.Error(err.Error())
				}

				break // use only stdin
			}

			select {
			case <-ctx.Done():
				return
			default:
				ldr.loadPath(ctx, heap.FromPath(sys.SplitPart(path)))
			}
		}
	}()

	return ldr.heaps
}

func (ldr *Loader) Exit() {
	slog.Info(fmt.Sprintf("total size %s", sys.Humanize(ldr.size.Load())))
}

func (ldr *Loader) loadPath(ctx context.Context, h *heap.Heap) {
	slog.Debug(fmt.Sprintf("looking for %s", h.String()))

	match, err := doublestar.FilepathGlob(h.Path)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	if len(match) == 0 {
		slog.Error(fmt.Sprintf("no files found for %s", h.String()))
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

			v := heap.FromPath(path, h.Part)

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

func (ldr *Loader) loadDir(ctx context.Context, h *heap.Heap, i int) error {
	if ldr.opts.Protect && i >= sys.MaxFolders {
		return errors.New("max folders reached")
	}

	dir, err := os.ReadDir(h.Path)

	if err != nil {
		return err
	}

	for _, f := range dir {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			v := heap.FromPath(filepath.Join(h.Path, f.Name()), h.Part)

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

func (ldr *Loader) loadFile(ctx context.Context, h *heap.Heap) error {
	f, err := os.Open(h.Path)

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
		return ldr.createHeap(ctx, h)
	}

	token, b, err := memory.Alloc(f)

	if err != nil {
		return err
	}

	h = heap.New(h.Path, h.Part, "", token, b)

	slog.Debug(fmt.Sprintf("memory mapped %s", h.String()))

	return ldr.processData(ctx, h, 0)
}

func (ldr *Loader) processData(ctx context.Context, h *heap.Heap, i int) error {
	// check depth to protect against zip bombs
	if ldr.opts.Protect && i >= sys.MaxArchives {
		return errors.New("max archives reached")
	}

	// stage 1: deflate data (nested)
	for j := 1; ; j++ {
		if ldr.opts.Protect && j >= sys.MaxDeflates {
			return errors.New("max deflates reached")
		}

		if !ldr.deflateData(h) {
			break // no more nested
		}
	}

	// stage 2: extract data (recursive)
	if ok, err := ldr.extractData(ctx, h, i); ok {
		return nil // ends here
	} else if err != nil {
		return err
	}

	// stage 3: convert data
	ldr.convertData(h)

	// stage 4: format data
	ldr.formatData(h)

	// filter for specific streams
	if len(h.Part) == 0 || strings.Contains(h.Path, h.Part) {
		return ldr.createHeap(ctx, h)
	}

	return nil
}

func (ldr *Loader) extractData(ctx context.Context, h *heap.Heap, i int) (bool, error) {
	if ldr.opts.Protect && ldr.files.Load() >= sys.MaxFiles {
		return false, errors.New("max files reached")
	}

	for _, r := range registry.Extracts {
		if r.Detect(h.Bytes()) {
			slog.Debug(fmt.Sprintf("archive detected as possibly %s", r.Name))

			for _, s := range r.Extract(h.Bytes(), h.Path, ldr.opts.Password) {
				slog.Debug(fmt.Sprintf("stream found %s", s.Path))

				select {
				case <-ctx.Done():
					return true, nil
				default:
					if err := ldr.processData(ctx, heap.New(
						s.Path,
						h.Part,
						"",
						0,
						s.Data,
					), i+1); err != nil {
						slog.Error(err.Error())
					}
				}
			}

			return true, nil
		}
	}

	return false, nil
}

func (ldr *Loader) deflateData(h *heap.Heap) bool {
	for _, r := range registry.Deflates {
		if r.Detect(h.Bytes()) {
			slog.Debug(fmt.Sprintf("deflate detected as possibly %s", r.Name))

			b, err := r.Deflate(h.Bytes())
			if err != nil {
				slog.Error(err.Error())
				return false
			}

			h.Change(b)
			return true
		}
	}

	return false
}

func (ldr *Loader) convertData(h *heap.Heap) {
	for _, r := range registry.Converts {
		if r.Detect(h.Bytes()) {
			slog.Debug(fmt.Sprintf("convert detected as possibly %s", r.Name))

			b, err := r.Convert(h.Bytes())
			if err != nil {
				slog.Error(err.Error())
				return
			}

			h.Change(b)
			h.Hint = "json"
			return
		}
	}
}

func (ldr *Loader) formatData(h *heap.Heap) {
	for _, r := range registry.Formats {
		if r.Detect(h.Bytes()) {
			slog.Debug(fmt.Sprintf("format detected as possibly %s", r.Name))

			b, err := r.Format(h.Bytes())
			if err != nil {
				slog.Error(err.Error())
				return
			}

			h.Change(b)
			h.Hint = "json"
			return
		}
	}
}

func (ldr *Loader) createHeap(ctx context.Context, h *heap.Heap) error {
	if ldr.opts.Protect && ldr.files.Load() >= sys.MaxFiles {
		return errors.New("max files reached")
	}

	if _, ok := ldr.paths.LoadOrStore(h.String(), pkg.Nil{}); ok {
		return nil // already loaded
	}

	ldr.size.Add(h.Size) // add original size
	ldr.files.Add(1)

	if b, ok := ldr.opts.Query.Reduce(h.Bytes()); ok {
		h.Change(b)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()

	case ldr.heaps <- h:
		slog.Debug(fmt.Sprintf("loaded heap %s", h.String()))
	}

	return nil
}
