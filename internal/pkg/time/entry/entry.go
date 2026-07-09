package entry

import (
	"fmt"
	"strings"
	"time"
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
	return e.AsBody()
}

func (e Entry) AsBody() string {
	return fmt.Sprintf("0|%s|%s|%s|0|0|%d|%d|%d|%d|%d",
		replacer.Replace(e.Name),
		e.Inode,
		e.Mode,
		e.Size,
		timeOrZero(e.Atime),
		timeOrZero(e.Mtime),
		timeOrZero(e.Ctime),
		timeOrZero(e.Btime),
	)
}

func (e Entry) AsTimesketch() string {
	var lines []string

	if !e.Mtime.IsZero() {
		lines = append(lines, fmt.Sprintf("%s,%d,%s,%s",
			e.Name+" was modified",
			e.Mtime.UTC().UnixMicro(),
			e.Mtime.UTC().Format(time.RFC3339),
			"Modify time",
		))
	}

	if !e.Atime.IsZero() {
		lines = append(lines, fmt.Sprintf("%s,%d,%s,%s",
			e.Name+" was accessed",
			e.Atime.UTC().UnixMicro(),
			e.Atime.UTC().Format(time.RFC3339),
			"Access time",
		))
	}

	if !e.Ctime.IsZero() {
		lines = append(lines, fmt.Sprintf("%s,%d,%s,%s",
			e.Name+" was changed",
			e.Ctime.UTC().UnixMicro(),
			e.Ctime.UTC().Format(time.RFC3339),
			"Change time",
		))
	}

	if !e.Btime.IsZero() {
		lines = append(lines, fmt.Sprintf("%s,%d,%s,%s",
			e.Name+" was created",
			e.Btime.UTC().UnixMicro(),
			e.Btime.UTC().Format(time.RFC3339),
			"Create time",
		))
	}

	return strings.Join(lines, "\n")
}

func timeOrZero(t time.Time) int64 {
	if !t.IsZero() {
		return t.UTC().Unix()
	}

	return 0
}
