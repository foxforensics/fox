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

const LocalHost = "http://localhost:8088/services/collector/event/1.0"

type Hec struct {
	Time       int64          `json:"time"`
	Host       string         `json:"host,omitempty"`
	Source     string         `json:"source"`
	Sourcetype string         `json:"sourcetype"`
	Index      string         `json:"index"`
	Fields     map[string]any `json:"fields,omitempty"`
	Event      Event          `json:"event"`

	// internal
	url, token string
}

type Event struct {
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

func New(url, token string) Hec {
	hec := Hec{url: url, token: token}

	hec.Index = "main"
	hec.Sourcetype = "_json"
	hec.Source = fmt.Sprintf("fox %s", res.Version)

	return hec
}

func (hec Hec) String() string {
	return fmt.Sprintf("HEC @ %s", hec.url)
}

func (hec Hec) Stream(e *event.Event) error {
	hec.Time = e.Time.UTC().UnixMilli()
	hec.Host = e.Host
	hec.Event = Event{
		e.Message,
		toCEF(e.Severity),
	}

	hec.Fields = make(map[string]any)

	for k, v := range e.Fields {
		hec.Fields[k] = v
	}

	buf, err := json.Marshal(hec)

	if err != nil {
		return err
	}

	return stream.Post(hec.url, string(buf), map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Splunk %s", strings.ToLower(hec.token)),
	})
}

func toCEF(n int) string {
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
