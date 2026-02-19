package help

import (
	"fmt"
	"strings"

	cli "foxhunt.dev/fox/internal/cmd"

	"foxhunt.dev/fox/internal/cmd/cat"
	"foxhunt.dev/fox/internal/cmd/dump"
	"foxhunt.dev/fox/internal/cmd/hash"
	"foxhunt.dev/fox/internal/cmd/hex"
	"foxhunt.dev/fox/internal/cmd/hunt"
	"foxhunt.dev/fox/internal/cmd/list"
	"foxhunt.dev/fox/internal/cmd/test"
	"foxhunt.dev/fox/internal/cmd/text"
)

var usage = map[string]string{
	"cat":  cat.Usage,
	"hex":  hex.Usage,
	"text": text.Usage,
	"hash": hash.Usage,
	"list": list.Usage,
	"dump": dump.Usage,
	"test": test.Usage,
	"hunt": hunt.Usage,
}

type Help struct {
	Mode string `arg:"" optional:""`
}

func (cmd *Help) Run(_ *cli.Globals) error {
	if v, ok := usage[strings.ToLower(cmd.Mode)]; ok {
		fmt.Println(v)
	}

	return nil
}
