package help

import (
	"strings"

	cli "go.foxforensics.dev/fox/v4/internal/cmd"

	"go.foxforensics.dev/fox/v4/internal/cmd/dump"
	"go.foxforensics.dev/fox/v4/internal/cmd/hash"
	"go.foxforensics.dev/fox/v4/internal/cmd/hunt"
	"go.foxforensics.dev/fox/v4/internal/cmd/info"
	"go.foxforensics.dev/fox/v4/internal/cmd/str"
	"go.foxforensics.dev/fox/v4/internal/pkg/text"
)

var usage = map[string]string{
	"str":  str.Usage,
	"info": info.Usage,
	"hash": hash.Usage,
	"dump": dump.Usage,
	"hunt": hunt.Usage,
}

type Help struct {
	Name string `arg:"" optional:""`
}

func (cmd *Help) Run(_ *cli.Globals) error {
	if v, ok := usage[strings.ToLower(cmd.Name)]; ok {
		return text.Usage(v)
	}

	return nil
}
