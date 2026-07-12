// Package hec applies this schema: https://docs.splunk.com/Documentation/Splunk/latest/Data/FormateventsforHTTPEventCollector
package hec

import (
	"encoding/json"
	"fmt"

	"go.foxforensics.eu/fox/v5/internal/cmd/hunt/event"
	"go.foxforensics.eu/fox/v5/internal/pkg/version"
)

type Hec struct {
	Time       int64          `json:"time"`
	Host       string         `json:"host,omitempty"`
	Source     string         `json:"source"`
	Sourcetype string         `json:"sourcetype"`
	Index      string         `json:"index"`
	Fields     map[string]any `json:"fields,omitempty"`
	Event      struct {
		Message  string `json:"message"`
		Severity string `json:"severity"`
	} `json:"event"`
}

func Apply(evt *event.Event) ([]byte, error) {
	hec := &Hec{
		Time:       evt.Time.UTC().UnixMilli(),
		Host:       evt.Host,
		Source:     fmt.Sprintf("fox %s", version.Number),
		Sourcetype: "_json",
		Index:      "main",
		Fields:     make(map[string]any),
	}

	for k, v := range evt.Fields {
		hec.Fields[k] = v
	}

	hec.Event.Message = evt.Message

	switch evt.Severity {
	case 10, 9:
		hec.Event.Severity = "CRITICAL"
	case 8, 7:
		hec.Event.Severity = "HIGH"
	case 6, 5, 4:
		hec.Event.Severity = "MEDIUM"
	default:
		hec.Event.Severity = "LOW"
	}

	return json.Marshal(hec)
}
