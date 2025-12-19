package raw

import (
	"fmt"

	"github.com/cuhsat/fox/v4/internal/pkg/data/stream"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
)

type Raw struct {
	stream.Stream
}

func New(url string) Raw {
	return Raw{stream.Stream{Url: url, Map: map[string]string{
		"Content-Type": "text/plain",
	}}}
}

func (raw Raw) String() string {
	return fmt.Sprintf("raw: %s", raw.Url)
}

func (raw Raw) Write(e *event.Event) (int64, int64, error) {
	return raw.Post(e.ToCEF())
}
