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
	"go.foxforensics.eu/fox/v4/internal/cmd/time"
	"go.foxforensics.eu/fox/v4/internal/sys"
)

var usage = map[string]string{
	"ad":   ad.Usage,
	"cat":  cat.Usage,
	"hash": hash.Usage,
	"hunt": hunt.Usage,
	"info": info.Usage,
	"str":  str.Usage,
	"time": time.Usage,
}

type Help struct {
	Name string `arg:"" optional:""`
}

func (cmd *Help) Run(_ *cmd.Globals) error {
	if v, ok := usage[strings.ToLower(cmd.Name)]; ok {
		sys.Usage(v)
		return nil
	}

	return errors.New("help topic is unknown")
}
