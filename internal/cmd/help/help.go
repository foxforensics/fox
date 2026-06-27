package help

import (
	"errors"
	"strings"

	"go.foxforensics.eu/fox/v4/internal/cmd"
	"go.foxforensics.eu/fox/v4/internal/cmd/ad"
	"go.foxforensics.eu/fox/v4/internal/cmd/cat"
	"go.foxforensics.eu/fox/v4/internal/cmd/hash"
	"go.foxforensics.eu/fox/v4/internal/cmd/hunt"
	"go.foxforensics.eu/fox/v4/internal/cmd/info"
	"go.foxforensics.eu/fox/v4/internal/cmd/str"
	"go.foxforensics.eu/fox/v4/internal/sys"
)

var usage = map[string]string{
	"ad":   ad.Usage,
	"cat":  cat.Usage,
	"str":  str.Usage,
	"info": info.Usage,
	"hash": hash.Usage,
	"hunt": hunt.Usage,
}

type Help struct {
	Name string `arg:"" optional:""`
}

func (cmd *Help) Run(_ *cmd.Globals) error {
	if v, ok := usage[strings.ToLower(cmd.Name)]; ok {
		return sys.Usage(v)
	}

	return errors.New("help topic is unknown")
}
