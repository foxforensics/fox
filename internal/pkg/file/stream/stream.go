package stream

import "go.foxforensics.dev/fox/v4/internal/pkg/types/event"

type Streamer interface {
	Stream(*event.Event) error
}
