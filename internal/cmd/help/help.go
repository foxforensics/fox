package help

import (
	"fmt"
	"strings"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/cmd/cat"
	"github.com/cuhsat/fox/v4/internal/cmd/hash"
	"github.com/cuhsat/fox/v4/internal/cmd/hex"
	"github.com/cuhsat/fox/v4/internal/cmd/hunt"
	"github.com/cuhsat/fox/v4/internal/cmd/list"
	"github.com/cuhsat/fox/v4/internal/cmd/test"
	"github.com/cuhsat/fox/v4/internal/cmd/text"
)

var usage = map[string]string{
	"cat":  cat.Usage,
	"hex":  hex.Usage,
	"list": list.Usage,
	"text": text.Usage,
	"test": test.Usage,
	"hash": hash.Usage,
	"hunt": hunt.Usage,
}

type Help struct {
	Mode string `arg:"" optional:""`
}

func (cmd *Help) Run(_ *cli.Globals) error {
	if v, ok := usage[strings.ToLower(cmd.Mode)]; ok {
		fmt.Print(v)
	}

	return nil
}
