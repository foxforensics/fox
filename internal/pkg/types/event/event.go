package event

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/cuhsat/fox/v4/internal"
)

const CEF = "%s %s CEF:1|fox|hunt|%s|100|%s|%d|"

type Event struct {
	Time      time.Time      `json:"ts"`
	Host      string         `json:"host,omitempty"`
	User      string         `json:"user,omitempty"`
	Message   string         `json:"msg,omitempty"`
	Severity  int8           `json:"lvl"`
	Extension map[string]any `json:"ext,omitempty"`
}

func (evt *Event) String() string {
	return evt.ToCEF()
}

func (evt *Event) ToCEF() string {
	var sb strings.Builder

	msg := evt.Message
	msg = strings.ReplaceAll(msg, `\`, `\\`)
	msg = strings.ReplaceAll(msg, `|`, `\|`)
	msg = strings.ReplaceAll(msg, `\t`, ` `)
	msg = strings.ReplaceAll(msg, `\n`, ``)

	if len(msg) > 512 {
		msg = msg[:512]
	}

	sb.WriteString(fmt.Sprintf(CEF,
		evt.Time.Format("Jan 02 2006 15:04:05.000"),
		evt.Host,
		app.Version[1:],
		msg,
		evt.Severity,
	))

	for _, k := range slices.Sorted(maps.Keys(evt.Extension)) {
		if v := evt.Extension[k]; v != nil {
			s := fmt.Sprintf("%v", v)

			k = strings.ReplaceAll(k, `=`, `\=`)
			s = strings.ReplaceAll(s, `=`, `\=`)

			sb.WriteString(fmt.Sprintf("%s=%s ", k, s))
		}
	}

	return strings.TrimSpace(sb.String())
}

func (evt *Event) ToJSON() string {
	b, _ := json.MarshalIndent(evt, "", "  ")

	return string(b)
}

func (evt *Event) ToJSONL() string {
	b, _ := json.Marshal(evt)

	return string(b)
}
