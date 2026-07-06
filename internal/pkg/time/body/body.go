package body

import (
	"fmt"
	"strings"
	"time"
)

var replacer = strings.NewReplacer(
	`\`, `\\`, // mask backslashes
	`|`, `\|`, // mask pipes
)

type Body struct {
	MD5    []byte    `json:"MD5,omitempty"`
	Name   string    `json:"name,omitempty"`
	Inode  string    `json:"inode,omitempty"`
	Mode   string    `json:"mode_as_string,omitempty"`
	UID    uint64    `json:"UID,omitempty"`
	GID    uint64    `json:"GID,omitempty"`
	Size   uint64    `json:"size,omitempty"`
	Atime  time.Time `json:"atime,omitempty"`
	Mtime  time.Time `json:"mtime,omitempty"`
	Ctime  time.Time `json:"ctime,omitempty"`
	Crtime time.Time `json:"crtime,omitempty"`
}

func (b Body) String() string {
	return fmt.Sprintf("0|%s|%s|%s|%d|%d|%d|%d|%d|%d|%d",
		replacer.Replace(b.Name),
		b.Inode,
		b.Mode,
		b.UID,
		b.GID,
		b.Size,
		b.Atime.UTC().Unix(),
		b.Mtime.UTC().Unix(),
		b.Ctime.UTC().Unix(),
		b.Crtime.UTC().Unix(),
	)
}
