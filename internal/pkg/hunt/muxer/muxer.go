package muxer

import (
	"context"
	"log/slog"

	"github.com/sourcegraph/conc/pool"
	"go.foxforensics.eu/fox/v5/internal/pkg/hunt/event"
)

type Handler func(context.Context, *event.Event) error

type Muxer struct {
	buffer   int
	handlers *pool.ContextPool
	channels []chan *event.Event
}

func New(ctx context.Context, n int) *Muxer {
	return &Muxer{
		buffer:   max(1, n),
		handlers: pool.New().WithContext(ctx),
		channels: make([]chan *event.Event, 0),
	}
}

func (m *Muxer) AddHandler(h Handler) {
	ch := make(chan *event.Event, m.buffer)

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

				if e == nil {
					continue // skipped
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
		return

	default:
		for _, ch := range m.channels {
			select {
			case ch <- e:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (m *Muxer) Wait() error {
	for _, ch := range m.channels {
		close(ch)
	}

	return m.handlers.Wait()
}
