package info

import (
	"fmt"
	"log"
	"strings"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
)

var Usage = strings.TrimSpace(`
Prints file infos and entropy.

fox info [FLAGS ...] <PATHS ...>

Flags:
  -m, --min=DECIMAL        minimum entropy value (default 0.0)
  -x, --max=DECIMAL        maximal entropy value (default 1.0)

Examples:
  $ fox info -m0.9 ./**/*
`)

type Info struct {
	Min   float64  `short:"m" default:"0.0"`
	Max   float64  `short:"x" default:"1.0"`
	Paths []string `arg:"" name:"path" type:"path" optional:""`
}

func (cmd *Info) Validate() error {
	if cmd.Min > cmd.Max {
		log.Fatal("invalid range")
	}

	return nil
}

func (cmd *Info) Run(cli *cli.Globals) error {
	if cli.Help || len(cmd.Paths) == 0 {
		fmt.Print(Usage)
		return nil
	}

	hs := cli.Bootstrap(cmd.Paths)
	defer cli.Discard()

	for _, h := range hs.Get() {
		if e, ok := h.Entropy(
			cmd.Min,
			cmd.Max,
		); ok {
			_, _ = fmt.Fprintf(cli.Stdout, "%10dL %10dB  %.10fE  %s\n", h.Len(), len(h.MMap()), e, text.Hide(h.String()))
		}
	}

	return nil
}
