package help

import (
	"strings"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/cmd/dump"
	"github.com/cuhsat/fox/v4/internal/cmd/hash"
	"github.com/cuhsat/fox/v4/internal/cmd/hunt"
	"github.com/cuhsat/fox/v4/internal/cmd/info"
	"github.com/cuhsat/fox/v4/internal/cmd/str"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
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
