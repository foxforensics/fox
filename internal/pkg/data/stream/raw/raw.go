package raw

import "github.com/cuhsat/fox/v4/internal/pkg/data/stream"

type Raw struct {
	stream.Schema
}

func New(url string) Raw {
	return Raw{stream.Schema{Url: url, Map: map[string]string{
		"Content-Type": "text/plain",
	}}}
}

func (raw Raw) Write(p []byte) (int, error) {
	return raw.Post(string(p))
}
