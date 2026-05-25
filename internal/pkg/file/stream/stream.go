package stream

import "go.foxforensics.dev/fox/v4/internal/pkg/types/event"

// local urls
const (
	Elastic = "http://localhost:8080"
	Splunk  = "http://localhost:8088/services/collector/event/1.0"
)

type Streamer interface {
	Stream(*event.Event) error
}
