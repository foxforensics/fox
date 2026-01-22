package info

import (
	"bytes"
	"fmt"
	"log"
	"slices"
	"strings"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heap"
)

var Usage = strings.TrimSpace(`
Prints infos and entropy.

fox info [FLAGS...] <PATHS...>

Flags:
  -b, --block=SIZE         block size for calculations
  -n, --min=DECIMAL        minimum entropy value (default: 0.0)
  -x, --max=DECIMAL        maximal entropy value (default: 1.0)

Example:
  $ fox info -n0.9 ./**/*
`)

type Info struct {
	Block uint64   `short:"b"`
	Min   float64  `short:"n" default:"0.0"`
	Max   float64  `short:"x" default:"1.0"`
	Paths []string `arg:"" name:"path" type:"path" optional:""`
}

func (cmd *Info) Validate() error {
	if cmd.Min > cmd.Max {
		log.Fatalln("invalid range")
	}

	return nil
}

func (cmd *Info) Run(cli *cli.Globals) error {
	if cli.Help || len(cmd.Paths)+len(cli.File) == 0 {
		fmt.Print(Usage)
		return nil
	}

	cli.NoConvert = true // forced

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	for h := range ch {
		var n = cmd.Block
		var off int

		if n == 0 {
			n = h.Size
		}

		if h.Size == 0 {
			_, _ = fmt.Fprintf(cli.Stdout, "%10dl %10db  %.10fe  %s\n", 0, 0, 0.0, text.Hide(h.String()))
			h.Discard()
			continue
		}

		for block := range slices.Chunk(h.Bytes(), int(n)) {
			l := bytes.Count(block, []byte{'\n'})
			e := heap.Entropy(block)

			// add possibly remaining line
			if block[len(block)-1] != '\n' {
				l++
			}

			if e >= cmd.Min && e <= cmd.Max {
				title := text.Hide(h.String())
				start := text.Hide(fmt.Sprintf("(@%d)", off))

				if cmd.Block > 0 {
					_, _ = fmt.Fprintf(cli.Stdout, "%10dl %10db  %.10fe  %s %s\n", l, len(block), e, title, start)
				} else {
					_, _ = fmt.Fprintf(cli.Stdout, "%10dl %10db  %.10fe  %s\n", l, len(block), e, title)
				}
			}

			off += len(block)
		}

		h.Discard()
	}

	return nil
}
