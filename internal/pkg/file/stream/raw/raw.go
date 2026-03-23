package raw

import (
	"fmt"

	"github.com/cuhsat/fox/v4/internal/pkg/file/stream"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
)

type Raw struct {
	url string
}

func New(url string) Raw {
	return Raw{url}
}

func (raw Raw) String() string {
	return fmt.Sprintf("raw @ %s", raw.url)
}

func (raw Raw) Stream(e *event.Event) error {
	return stream.Post(raw.url, e.ToCEF(), map[string]string{
		"Content-Type": "text/plain",
	})
}
