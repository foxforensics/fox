package hunter

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"strings"
	"sync"

	"github.com/sourcegraph/conc/pool"
	"go.foxforensics.eu/fox/v4/internal/lib/binaries/log/evtx"
	"go.foxforensics.eu/fox/v4/internal/lib/binaries/log/journal"
	"go.foxforensics.eu/fox/v4/internal/pkg/hunt/event"
	"go.foxforensics.eu/fox/v4/internal/sys/heap"
)

const (
	Latency = int64(1024 * 64) // 64kb
	Scale   = 1024
)

var Local = []string{
	"/Windows/System32/winevt/Logs",
	"/var/log/journal",
	"/run/log/journal",
}

type Options struct {
	Sort    bool
	Threads int
}

type Hunter struct {
	opts *Options
}

func New(opts *Options) *Hunter {
	sync.OnceFunc(evtx.Prepare)()

	return &Hunter{
		opts: opts,
	}
}

func (htr *Hunter) Hunt(ctx context.Context, heaps <-chan *heap.Heap) <-chan *event.Event {
	ch := make(chan *event.Event, htr.opts.Threads*Scale)

	go func() {
		defer close(ch)

		p := pool.New().
			WithContext(ctx).
			WithMaxGoroutines(htr.opts.Threads)

		for h := range heaps {
			p.Go(func(ctx context.Context) error {
				slog.Info(fmt.Sprintf("hunt: carving %s", h.String()))

				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					return htr.carve(ctx, ch, h)
				}
			})
		}

		if err := p.Wait(); err != nil {
			if errors.Is(err, context.Canceled) {
				slog.Info("hunt: canceled")
			} else {
				slog.Error(err.Error())
			}
		}
	}()

	if htr.opts.Sort {
		return htr.sort(ch)
	}

	return ch
}

func (htr *Hunter) sort(ch <-chan *event.Event) <-chan *event.Event {
	sorted := make(chan *event.Event, cap(ch))

	go func() {
		defer close(sorted)

		var src []*event.Event

		for e := range ch {
			src = append(src, e)
		}

		slices.SortStableFunc(src, func(a, b *event.Event) int {
			return strings.Compare(a.SortKey(), b.SortKey())
		})

		for _, e := range src {
			sorted <- e
		}
	}()

	return sorted
}

func (htr *Hunter) carve(ctx context.Context, ch chan<- *event.Event, h *heap.Heap) error {
	defer h.Discard()

	p := pool.New().
		WithContext(ctx).
		WithMaxGoroutines(htr.opts.Threads)

	p.Go(func(ctx context.Context) error {
		return htr.carveEvtx(ctx, ch, h)
	})

	p.Go(func(ctx context.Context) error {
		return htr.carveJournal(ctx, ch, h)
	})

	return p.Wait()
}

func (htr *Hunter) carveEvtx(ctx context.Context, ch chan<- *event.Event, h *heap.Heap) error {
	sr := io.NewSectionReader(h.Reader(), 0, int64(h.Size))
	for off := range htr.findOffset(ctx, h, evtx.Chunk) {
		for evt := range evtx.Carve(sr, off, cap(ch)) {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				ch <- evt
			}
		}
	}

	return nil
}

func (htr *Hunter) carveJournal(ctx context.Context, ch chan<- *event.Event, h *heap.Heap) error {
	sr := io.NewSectionReader(h.Reader(), 0, int64(h.Size))
	for off := range htr.findOffset(ctx, h, journal.Magic) {
		for evt := range journal.Carve(sr, off, cap(ch)) {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				ch <- evt
			}
		}
	}

	return nil
}

func (htr *Hunter) findOffset(ctx context.Context, h *heap.Heap, seq []byte) <-chan int64 {
	out := make(chan int64, 64*htr.opts.Threads)

	go func(b []byte) {
		defer close(out)

		var off, idx int64
		for off < int64(len(b)) {
			if idx = int64(bytes.Index(b[off:], seq)); idx == -1 {
				break
			}

			off += idx
			out <- off

			slog.Debug(fmt.Sprintf("hunt: found at offset 0x%08x", off))

			off += int64(len(seq))

			// we trade off carving speed over latency
			if off%Latency == 0 {
				select {
				case <-ctx.Done():
					return // hunt canceled
				default:
					continue
				}
			}
		}
	}(h.Bytes())

	return out
}
