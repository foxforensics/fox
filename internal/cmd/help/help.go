package help

import (
	"strings"

	cli "github.com/cuhsat/fox/v4/internal/cmd"
	std "github.com/cuhsat/fox/v4/internal/pkg/text"

	"github.com/cuhsat/fox/v4/internal/cmd/cat"
	"github.com/cuhsat/fox/v4/internal/cmd/dump"
	"github.com/cuhsat/fox/v4/internal/cmd/hash"
	"github.com/cuhsat/fox/v4/internal/cmd/hex"
	"github.com/cuhsat/fox/v4/internal/cmd/hunt"
	"github.com/cuhsat/fox/v4/internal/cmd/mcp"
	"github.com/cuhsat/fox/v4/internal/cmd/stat"
	"github.com/cuhsat/fox/v4/internal/cmd/test"
	"github.com/cuhsat/fox/v4/internal/cmd/text"
)

var usage = map[string]string{
	"mcp":  mcp.Usage,
	"cat":  cat.Usage,
	"hex":  hex.Usage,
	"text": text.Usage,
	"hash": hash.Usage,
	"stat": stat.Usage,
	"dump": dump.Usage,
	"test": test.Usage,
	"hunt": hunt.Usage,
}

type Help struct {
	Mode string `arg:"" optional:""`
}

func (cmd *Help) Run(_ *cli.Globals) error {
	if v, ok := usage[strings.ToLower(cmd.Mode)]; ok {
		return std.Usage(v)
	}

	return nil
}
