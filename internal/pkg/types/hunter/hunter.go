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

	"go.foxforensics.eu/fox/v4/internal/pkg/file/binary/log/evtx"
	"go.foxforensics.eu/fox/v4/internal/pkg/file/binary/log/journal"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/event"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/heap"
)

var Latency = int64(1024 * 1024) // 1mb

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
	opts   *Options
	events chan *event.Event
}

func New(opts *Options) *Hunter {
	return &Hunter{
		opts:   opts,
		events: make(chan *event.Event, opts.Threads*1024),
	}
}

func (htr *Hunter) Hunt(ctx context.Context, ch <-chan *heap.Heap) <-chan *event.Event {
	go func() {
		defer close(htr.events)

		p := pool.New().
			WithContext(ctx).
			WithFirstError().
			WithMaxGoroutines(htr.opts.Threads)

		for h := range ch {
			p.Go(func(ctx context.Context) error {
				slog.Info(fmt.Sprintf("hunt: carving %s", h.String()))

				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					return htr.carve(ctx, h)
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
		return htr.sort()
	}

	return htr.events
}

func (htr *Hunter) sort() <-chan *event.Event {
	sorted := make(chan *event.Event, cap(htr.events))

	go func() {
		defer close(sorted)

		var src []*event.Event

		for e := range htr.events {
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

func (htr *Hunter) carve(ctx context.Context, h *heap.Heap) error {
	defer h.Discard()

	sync.OnceFunc(evtx.Preload)()

	p := pool.New().
		WithContext(ctx).
		WithMaxGoroutines(htr.opts.Threads)

	p.Go(func(ctx context.Context) error {
		return htr.carveEvtx(ctx, h)
	})

	p.Go(func(ctx context.Context) error {
		return htr.carveJournal(ctx, h)
	})

	return p.Wait()
}

func (htr *Hunter) carveEvtx(ctx context.Context, h *heap.Heap) error {
	sr := io.NewSectionReader(h.Reader(), 0, int64(h.Size))
	for off := range htr.findOffset(ctx, h, evtx.Chunk) {
		for evt := range evtx.Carve(sr, off, cap(htr.events)) {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				htr.events <- evt
			}
		}
	}

	return nil
}

func (htr *Hunter) carveJournal(ctx context.Context, h *heap.Heap) error {
	sr := io.NewSectionReader(h.Reader(), 0, int64(h.Size))
	for off := range htr.findOffset(ctx, h, journal.Magic) {
		for evt := range journal.Carve(sr, off, cap(htr.events)) {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				htr.events <- evt
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
