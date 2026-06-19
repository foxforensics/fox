package stream

import (
	"context"

	"go.foxforensics.eu/fox/v4/internal/cmd/hunt/event"
)

// local urls
const (
	Elastic = "http://localhost:8080"
	Splunk  = "http://localhost:8088/services/collector/event/1.0"
)

type Streamer interface {
	Stream(context.Context, *event.Event) error
	Close() error
}
