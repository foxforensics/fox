// Package hec specification:
// https://docs.splunk.com/Documentation/Splunk/latest/Data/FormateventsforHTTPEventCollector
package hec

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cuhsat/fox/v4/internal"
	"github.com/cuhsat/fox/v4/internal/pkg/data/stream"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
)

type Hec struct {
	stream.Stream

	Time   int64  `json:"time"`
	Source string `json:"source"`
	Event  string `json:"event"`
}

func New(url, token string) Hec {
	return Hec{
		Source: fmt.Sprintf("fox %s", app.Version),
		Stream: stream.Stream{Url: url, Map: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Splunk %s", strings.ToLower(token)),
		}},
	}
}

func (hec Hec) String() string {
	return fmt.Sprintf("HEC: %s", hec.Url)
}

func (hec Hec) Write(e *event.Event) error {
	hec.Time = time.Now().UTC().UnixMilli()
	hec.Event = strings.TrimRight(e.ToCEF(), "\n")

	buf, err := json.Marshal(hec)

	if err != nil {
		return nil
	}

	return hec.Post(string(buf))
}
