package raw

import (
	"fmt"

	"github.com/cuhsat/fox/v4/internal/pkg/data/stream"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
)

type Raw struct {
	stream.Streamer
}

func New(url string) Raw {
	return Raw{stream.Streamer{
		Url: url,
		Map: map[string]string{
			"Content-Type": "text/plain",
		},
	}}
}

func (raw Raw) String() string {
	return fmt.Sprintf("raw @ %s", raw.Url)
}

func (raw Raw) Write(e *event.Event) error {
	return raw.Post(e.ToCEF())
}
