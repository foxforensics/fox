// Package hec specification:
// https://docs.splunk.com/Documentation/Splunk/latest/Data/FormateventsforHTTPEventCollector
package hec

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cuhsat/fox/v4/internal"
	"github.com/cuhsat/fox/v4/internal/pkg/data/stream"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
)

const Index = "main"
const Type = "_json"

type Hec struct {
	stream.Stream

	Time       int64          `json:"time"`
	Host       string         `json:"host,omitempty"`
	Source     string         `json:"source"`
	Sourcetype string         `json:"sourcetype"`
	Index      string         `json:"index"`
	Fields     map[string]any `json:"fields,omitempty"`
	Event      Event          `json:"event"`
}

type Event struct {
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

func New(url, token string) Hec {
	return Hec{
		Index:      Index,
		Sourcetype: Type,
		Source:     fmt.Sprintf("fox %s", app.Version[1:]),
		Stream: stream.Stream{Url: url, Map: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Splunk %s", strings.ToLower(token)),
		}},
	}
}

func (hec Hec) String() string {
	return fmt.Sprintf("HEC: %s", hec.Url)
}

func (hec Hec) Write(e *event.Event) (int64, int64, error) {
	hec.Time = e.Time.UTC().UnixMilli()
	hec.Host = e.Host
	hec.Event = Event{
		e.Message,
		cefName(e.Severity),
	}

	hec.Fields = make(map[string]any)

	for k, v := range e.Extension {
		hec.Fields[k] = v
	}

	buf, err := json.Marshal(hec)

	if err != nil {
		return 0, 0, nil
	}

	return hec.Post(string(buf))
}

func cefName(n int8) string {
	switch n {
	case 10, 9:
		return "CRITICAL"
	case 8, 7:
		return "HIGH"
	case 6, 5, 4:
		return "MEDIUM"
	default:
		return "LOW"
	}
}
