package list

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
Prints file infos and entropy.

fox list [FLAGS...] <PATHS...>

Flags:
  -b, --block=SIZE         block size for calculations
  -n, --min=VALUE          minimum entropy value (default: 0.0)
  -x, --max=VALUE          maximal entropy value (default: 1.0)

Format Flags:
  -H, --human              format size in human-readable units

Examples:
  $ fox list -n0.9 ./**/*
`)

type List struct {
	Block uint64  `short:"b"`
	Min   float64 `short:"n" default:"0.0"`
	Max   float64 `short:"x" default:"1.0"`

	// format
	Human bool `short:"H"`

	// paths
	Paths []string `arg:"" name:"path" type:"path" optional:""`
}

func (cmd *List) Validate() error {
	if cmd.Min > cmd.Max {
		log.Fatalln("invalid range")
	}

	return nil
}

func (cmd *List) Run(cli *cli.Globals) error {
	if len(cmd.Paths)+len(cli.Paths) == 0 {
		fmt.Println(Usage)
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
				size := fmt.Sprintf("%db", len(block))
				title := text.Hide(h.String())
				start := text.Hide(fmt.Sprintf("@ 0x%x", off))

				if cmd.Human {
					size = text.Humanize(int64(len(block)))
				}

				if cmd.Block > 0 {
					_, _ = fmt.Fprintf(cli.Stdout, "%10dl %10s  %.10fe  %s %s\n", l, size, e, title, start)
				} else {
					_, _ = fmt.Fprintf(cli.Stdout, "%10dl %10s  %.10fe  %s\n", l, size, e, title)
				}
			}

			off += len(block)
		}

		h.Discard()
	}

	return nil
}
