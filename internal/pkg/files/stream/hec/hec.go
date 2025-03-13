// Package hec specification:
// https://docs.splunk.com/Documentation/Splunk/latest/Data/FormateventsforHTTPEventCollector
package hec

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cuhsat/fox/v4/internal"
	"github.com/cuhsat/fox/v4/internal/pkg/files/stream"
)

type Hec struct {
	stream.Schema

	Time   int64  `json:"time"`
	Source string `json:"source"`
	Event  string `json:"event"`
}

func New(url, token string) Hec {
	return Hec{
		Source: fmt.Sprintf("fox %s", app.Version),
		Schema: stream.Schema{Url: url, Map: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Splunk %s", strings.ToLower(token)),
		}},
	}
}

func (hec Hec) Write(p []byte) (int, error) {
	hec.Time = time.Now().UTC().UnixMilli()
	hec.Event = strings.TrimRight(string(p), "\n")

	buf, err := json.Marshal(hec)

	if err != nil {
		return 0, nil
	}

	return hec.Post(string(buf))
}
