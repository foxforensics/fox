package entry

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"go.foxforensics.eu/fox/v5/library/formats"
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

type Timesketch struct {
	Message       string `json:"message"`
	Datetime      string `json:"datetime"`
	TimestampDesc string `json:"timestamp_desc"`
}

func (e Entry) String() string {
	return fmt.Sprintf("%s %s", e.SortKey().Format(time.RFC3339), e.Name)
}

func (e Entry) SortKey() time.Time {
	if !e.Ctime.IsZero() {
		return e.Ctime
	}

	if !e.Atime.IsZero() {
		return e.Atime
	}

	if !e.Mtime.IsZero() {
		return e.Mtime
	}

	if !e.Btime.IsZero() {
		return e.Btime
	}

	slog.Error("entry has no timestamp")
	return time.Time{}
}

func (e Entry) AsBodyfile() string {
	return fmt.Sprintf("0|%s|%s|%s|0|0|%d|%d|%d|%d|%d",
		replacer.Replace(e.Name),
		e.Inode,
		e.Mode,
		e.Size,
		toEpoch(e.Atime),
		toEpoch(e.Mtime),
		toEpoch(e.Ctime),
		toEpoch(e.Btime),
	)
}

func (e Entry) AsTimesketch() string {
	var sb strings.Builder

	sb.WriteString(toJsonl(e.Mtime, e.Name+" was modified", "Modify time"))
	sb.WriteString(toJsonl(e.Atime, e.Name+" was accessed", "Access time"))
	sb.WriteString(toJsonl(e.Ctime, e.Name+" was changed", "Change time"))
	sb.WriteString(toJsonl(e.Btime, e.Name+" was created", "Create time"))

	return sb.String()
}

func toJsonl(t time.Time, msg, desc string) string {
	if t.IsZero() {
		return ""
	}

	return formats.AsJSONL(&Timesketch{
		Message:       msg,
		Datetime:      t.Format(time.RFC3339),
		TimestampDesc: desc,
	})
}

func toEpoch(t time.Time) int64 {
	if !t.IsZero() {
		return t.Unix()
	}

	return 0
}
