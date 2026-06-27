package muxer

import (
	"context"
	"log/slog"

	"github.com/sourcegraph/conc/pool"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/event"
)

type Handler func(context.Context, *event.Event) error

type Muxer struct {
	scale    int
	threads  int
	handlers *pool.ContextPool
	channels []chan *event.Event
}

func New(ctx context.Context, t, n int) *Muxer {
	return &Muxer{
		scale:   max(1, n),
		threads: max(1, t),
		handlers: pool.New().
			WithContext(ctx).
			WithMaxGoroutines(t),
		channels: make([]chan *event.Event, 0, t),
	}
}

func (m *Muxer) AddHandler(h Handler) {
	ch := make(chan *event.Event, m.threads*m.scale)

	m.channels = append(m.channels, ch)
	m.handlers.Go(func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()

			case e, ok := <-ch:
				if !ok {
					return nil // closed
				}

				if err := h(ctx, e); err != nil {
					slog.Error(err.Error())
				}
			}
		}
	})
}

func (m *Muxer) Handle(ctx context.Context, e *event.Event) {
	select {
	case <-ctx.Done():
	default:
		for _, ch := range m.channels {
			ch <- e
		}
	}
}

func (m *Muxer) Wait() error {
	for _, ch := range m.channels {
		close(ch)
	}

	return m.handlers.Wait()
}
