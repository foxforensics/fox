package help

import (
	"strings"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/cmd/cat"
	"github.com/cuhsat/fox/v4/internal/cmd/check"
	"github.com/cuhsat/fox/v4/internal/cmd/dump"
	"github.com/cuhsat/fox/v4/internal/cmd/hash"
	"github.com/cuhsat/fox/v4/internal/cmd/hex"
	"github.com/cuhsat/fox/v4/internal/cmd/hunt"
	"github.com/cuhsat/fox/v4/internal/cmd/mcp"
	"github.com/cuhsat/fox/v4/internal/cmd/stat"
	"github.com/cuhsat/fox/v4/internal/cmd/str"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
)

var usage = map[string]string{
	"cat":   cat.Usage,
	"hex":   hex.Usage,
	"str":   str.Usage,
	"hash":  hash.Usage,
	"stat":  stat.Usage,
	"check": check.Usage,
	"dump":  dump.Usage,
	"hunt":  hunt.Usage,
	"mcp":   mcp.Usage,
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
