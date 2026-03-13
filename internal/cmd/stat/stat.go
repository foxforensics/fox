package stat

import (
	"bytes"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/alecthomas/kong"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heap"
)

var Usage = strings.TrimSpace(`
Show file stats and entropy.

fox stat [FLAGS...] <PATHS...>

Flags:
  -s, --sort               Sort files by path (slower)
  -b, --block=SIZE         Block size for analysis

Filter flags:
  -n, --min=VALUE          Minimum entropy value (default: 0.0)
  -x, --max=VALUE          Maximal entropy value (default: 1.0)

Format flags:
  -H, --human              Format size in human-readable units

Examples:
  $ fox stat -n0.8 ./**/*
`)

const Limit = 0.72

type Stat struct {
	Sort  bool    `short:"s"`
	Block string  `short:"b"`
	Min   float64 `short:"n" default:"0.0"`
	Max   float64 `short:"x" default:"1.0"`

	// format
	Human bool `short:"H"`

	// paths
	Paths []string `arg:"" name:"path" optional:""`

	// internal
	block uint64 `kong:"-"`
}

func (cmd *Stat) Validate() error {
	if cmd.Min > cmd.Max {
		log.Fatalln("invalid range")
	}

	return nil
}

func (cmd *Stat) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	if len(cmd.Block) > 0 {
		cmd.block = uint64(text.Mechanize(cmd.Block))
	}

	return nil
}

func (cmd *Stat) Run(cli *cli.Globals) error {
	if len(cmd.Paths)+len(cli.Paths) == 0 {
		return text.Usage(Usage)
	}

	if cmd.Sort {
		cli.Threads = 1 // single threaded
	}

	cli.NoConvert = true // forced

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	for h := range ch {
		var t = time.UnixMilli(int64(h.Time)).UTC().Format(time.RFC3339)
		var n = cmd.block
		var off int

		if n == 0 {
			n = h.Size
		}

		if h.Size == 0 {
			title := h.String()

			if cmd.block > 0 {
				title = "[00000000] " + title
			}

			text.Writeln(text.AsGray(fmt.Sprintf("%10dl %10db  %.10fe  %s  %s", 0, 0, 0.0, t, title)))

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
				title := text.AsBold(h.String())
				start := fmt.Sprintf("[%08x]", off)
				entropy := fmt.Sprintf("%.10fe", e)

				if cmd.Human {
					size = text.Humanize(int64(len(block)))
				}

				if e >= Limit {
					entropy = text.AsWarn(entropy)
				}

				if cmd.block > 0 {
					text.Writeln("%10dl %11s  %s  %s  %s %s", l, size, entropy, t, start, title)
				} else {
					text.Writeln("%10dl %11s  %s  %s  %s", l, size, entropy, t, title)
				}
			}

			off += len(block)
		}

		h.Discard()
	}

	return nil
}
