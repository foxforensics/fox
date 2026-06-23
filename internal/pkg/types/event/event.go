package event

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/zeebo/xxh3"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/fox/v4/internal/sys/version"
)

var replacer = strings.NewReplacer(
	"\t", " ", // remove tabs
	"\n", " ", // remove line breaks
	"\r", " ", // remove line breaks
	`\`, `\\`, // mask backslashes
	`|`, `\|`, // mask pipes
	`=`, `\=`, // mask equal
)

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

func (e *Event) AsCEF() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s %s CEF:1|fox|hunt|%s|100|",
		e.Time.Format("Jan 02 2006 15:04:05.000"),
		e.Host,
		version.Number,
	))

	_, _ = replacer.WriteString(&sb, sys.Sanitize(e.Message[:min(len(e.Message), 512)]))

	sb.WriteString(fmt.Sprintf("|%d|", e.Severity))

	for _, k := range slices.Sorted(maps.Keys(e.Fields)) {
		if v := e.Fields[k]; len(v) > 0 {
			sb.WriteString(fmt.Sprintf("%s=%s ",
				replacer.Replace(k),
				replacer.Replace(v),
			))
		}
	}

	return strings.TrimSpace(sb.String())
}
