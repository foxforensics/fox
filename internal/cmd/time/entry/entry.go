package entry

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"go.foxforensics.eu/fox/v5/internal/pkg/writer"
)

var replacer = strings.NewReplacer(
	`\`, `\\`, // mask backslashes
	`|`, `\|`, // mask pipes
	`/`, `\/`, // mask slashes
	`:`, `\:`, // mask colons
)

type Entry struct {
	Name    string    `json:"name,omitempty"`
	Inode   string    `json:"inode,omitempty"`
	Size    uint64    `json:"size"`
	Mode    string    `json:"mode,omitempty"`
	Mtime   time.Time `json:"m_time,omitempty"`
	Atime   time.Time `json:"a_time,omitempty"`
	Ctime   time.Time `json:"c_time,omitempty"`
	Btime   time.Time `json:"b_time,omitempty"`
	Anomaly bool      `json:"anomaly,omitempty"`
}

func (e Entry) String() string {
	s := fmt.Sprintf("0|%s|%s|%s|0|0|%d|%d|%d|%d|%d",
		replacer.Replace(e.Name),
		e.Inode,
		e.Mode,
		e.Size,
		epoch(e.Atime),
		epoch(e.Mtime),
		epoch(e.Ctime),
		epoch(e.Btime),
	)

	if e.Anomaly {
		return writer.AsBold(s)
	}

	return s
}

func (e Entry) SortKey() time.Time {
	if !e.Atime.IsZero() {
		return e.Atime
	}

	if !e.Mtime.IsZero() {
		return e.Mtime
	}

	if !e.Ctime.IsZero() {
		return e.Ctime
	}

	if !e.Btime.IsZero() {
		return e.Btime
	}

	slog.Error("entry has no timestamp")
	return time.Time{}
}

func epoch(t time.Time) int64 {
	if !t.IsZero() {
		return t.Unix()
	}

	return 0
}
