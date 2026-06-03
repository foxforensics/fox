package event

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/zeebo/xxh3"

	"go.foxforensics.dev/fox/v4/internal/pkg/text"
	"go.foxforensics.dev/fox/v4/internal/pkg/version"
)

const CEF = "%s %s CEF:1|fox|hunt|%s|100|%s|%d|"

type Event struct {
	Time     time.Time         `json:"time,omitempty"`
	Host     string            `json:"host,omitempty"`
	User     string            `json:"user,omitempty"`
	Message  string            `json:"message,omitempty"`
	Severity int               `json:"severity,omitempty"`
	Sequence string            `json:"sequence,omitempty"`
	Source   string            `json:"source,omitempty"`
	Category string            `json:"category,omitempty"`
	Service  string            `json:"service,omitempty"`
	Fields   map[string]string `json:"fields,omitempty"`
}

func (e *Event) String() string {
	return fmt.Sprintf("%s:%s:%s",
		e.Host,
		e.Message,
		e.Sequence,
	)
}

func (e *Event) SortKey() string {
	return fmt.Sprintf("%d-%d", e.Time.UnixNano(), xxh3.HashString(e.String()))
}

func (e *Event) ToCEF() string {
	var sb strings.Builder

	msg := text.Sanitize(e.Message)
	msg = strings.ReplaceAll(msg, `\`, `\\`)
	msg = strings.ReplaceAll(msg, `|`, `\|`)
	msg = strings.ReplaceAll(msg, `\t`, ` `)
	msg = strings.ReplaceAll(msg, `\n`, ``)

	if len(msg) > 512 {
		msg = msg[:512]
	}

	sb.WriteString(fmt.Sprintf(CEF,
		e.Time.Format("Jan 02 2006 15:04:05.000"),
		e.Host,
		version.Number,
		msg,
		e.Severity,
	))

	for _, k := range slices.Sorted(maps.Keys(e.Fields)) {
		if v := e.Fields[k]; len(v) > 0 {
			k = strings.ReplaceAll(k, `=`, `\=`)
			v = strings.ReplaceAll(v, `=`, `\=`)
			v = strings.ReplaceAll(v, "\n", " ")
			v = strings.ReplaceAll(v, "\r", "")
			v = strings.ReplaceAll(v, "\t", "")

			sb.WriteString(fmt.Sprintf("%s=%s ", k, v))
		}
	}

	return strings.TrimSpace(sb.String())
}

func (e *Event) ToJSON() string {
	b, _ := json.MarshalIndent(e, "", "  ")
	return string(b)
}

func (e *Event) ToJSONL() string {
	b, _ := json.Marshal(e)
	return string(b)
}
