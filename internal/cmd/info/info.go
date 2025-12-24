package info

import (
	"bytes"
	"fmt"
	"log"
	"slices"
	"strings"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
)

var Usage = strings.TrimSpace(`
Prints file infos and entropy.

fox info [FLAGS ...] <PATHS ...>

Flags:
  -b, --block=SIZE         block size for calculations
  -m, --min=DECIMAL        minimum entropy value (default 0.0)
  -x, --max=DECIMAL        maximal entropy value (default 1.0)

Examples:
  $ fox info -m0.9 ./**/*
`)

type Info struct {
	Block int64    `short:"b"`
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

	hs := cli.Load(cmd.Paths)
	defer cli.Discard()

	for _, h := range hs.Get() {
		var n = cmd.Block
		var off int

		if n == 0 {
			n = h.Size()
		}

		if h.Size() == 0 {
			_, _ = fmt.Fprintf(cli.Stdout, "%10dl %10db  %.10fe  %s\n", 0, 0, 0.0, text.Hide(h.String()))
			continue
		}

		for block := range slices.Chunk(h.MMap(), int(n)) {
			l := bytes.Count(block, []byte{'\n'})
			b := len(block)
			e := h.Entropy(block)

			// add possibly remaining line
			if block[len(block)-1] != '\n' {
				l++
			}

			if e >= cmd.Min && e <= cmd.Max {
				title := text.Hide(h.String())
				start := text.Hide(fmt.Sprintf("@%d", off))

				if cmd.Block > 0 {
					_, _ = fmt.Fprintf(cli.Stdout, "%10dl %10db  %.10fe  %s %s\n", l, b, e, title, start)
				} else {
					_, _ = fmt.Fprintf(cli.Stdout, "%10dl %10db  %.10fe  %s\n", l, b, e, title)
				}
			}

			off += b
		}
	}

	return nil
}
