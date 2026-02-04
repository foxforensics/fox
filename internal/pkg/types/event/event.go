package event

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/zeebo/xxh3"

	"github.com/cuhsat/fox/v4/internal"
)

const CEF = "%s %s CEF:1|fox|hunt|%s|100|%s|%d|"

type Event struct {
	Time     time.Time         `json:"time,omitempty"`
	Host     string            `json:"host,omitempty"`
	User     string            `json:"user,omitempty"`
	Message  string            `json:"message,omitempty"`
	Severity int8              `json:"severity,omitempty"`
	Sequence string            `json:"sequence,omitempty"`
	Source   string            `json:"source,omitempty"`
	Category string            `json:"category,omitempty"`
	Service  string            `json:"service,omitempty"`
	Fields   map[string]string `json:"fields,omitempty"`
}

func (evt *Event) String() string {
	return fmt.Sprintf("%s:%s:%s",
		evt.Host,
		evt.Message,
		evt.Sequence,
	)
}

func (evt *Event) SortKey() string {
	return fmt.Sprintf("%d-%d", evt.Time.UnixNano(), xxh3.HashString(evt.String()))
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
		app.Version,
		msg,
		evt.Severity,
	))

	ext := map[string]any{
		"rt":         evt.Time,
		"app":        evt.Source,
		"cat":        evt.Category,
		"sproc":      evt.Service,
		"shost":      evt.Host,
		"suid":       evt.User,
		"externalId": evt.Sequence,
	}

	for _, k := range slices.Sorted(maps.Keys(ext)) {
		if v := ext[k]; v != nil {
			if s := fmt.Sprintf("%v", v); len(s) > 0 {
				k = strings.ReplaceAll(k, `=`, `\=`)
				s = strings.ReplaceAll(s, `=`, `\=`)

				sb.WriteString(fmt.Sprintf("%s=%s ", k, s))
			}
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
