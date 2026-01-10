package info

import (
	"bytes"
	"fmt"
	"log"
	"slices"
	"strings"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/data/extern/vt"
	"github.com/cuhsat/fox/v4/internal/pkg/hash"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heap"
)

var Usage = strings.TrimSpace(`
Prints file infos and entropy.

fox info [FLAGS...] <PATHS...>

Flags:
  -b, --block=SIZE         block size for calculations
  -m, --min=DECIMAL        minimum entropy value (default: 0.0)
  -x, --max=DECIMAL        maximal entropy value (default: 1.0)

Check:
  -V, --vt=[=LEVEL]        check file using VirusTotal (V/VV/VVV)
  -K, --vt-key=KEY         check file using VirusTotal key

Example:
  $ fox info -m0.9 ./**/*
`)

type Info struct {
	Block int64    `short:"b"`
	Min   float64  `short:"m" default:"0.0"`
	Max   float64  `short:"x" default:"1.0"`
	Vt    int      `short:"V" long:"vt" type:"counter"`
	VtKey string   `short:"K" long:"vt-key"`
	Paths []string `arg:"" name:"path" type:"path" optional:""`
}

func (cmd *Info) Validate() error {
	if cmd.Min > cmd.Max {
		log.Fatalln("invalid range")
	}

	if cmd.Vt > 0 && len(cmd.VtKey) == 0 {
		log.Fatalln("api key required")
	}

	return nil
}

func (cmd *Info) Run(cli *cli.Globals) error {
	if cli.Help || len(cmd.Paths) == 0 {
		fmt.Print(Usage)
		return nil
	}

	cli.NoConvert = true // forced

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	for h := range ch {
		var n = cmd.Block
		var off int

		if cmd.Vt > 0 {
			var res string
			var err error

			if !cli.NoFile {
				_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide(text.Header(h.String())))
			}

			sum := hash.MustSum(types.SHA256, h.MMap())

			switch {
			case cmd.Vt > 2:
				res, err = vt.GetReport(sum, cmd.VtKey, cli.NoPretty)
			case cmd.Vt > 1:
				res, err = vt.GetResult(sum, cmd.VtKey)
			default:
				res, err = vt.GetLabel(sum, cmd.VtKey)
			}

			if err == nil {
				_, _ = fmt.Fprintln(cli.Stdout, res)
			} else {
				log.Println(err)
			}

			h.Discard()
			continue
		}

		if n == 0 {
			n = h.Size
		}

		if h.Size == 0 {
			_, _ = fmt.Fprintf(cli.Stdout, "%10dl %10db  %.10fe  %s\n", 0, 0, 0.0, text.Hide(h.String()))
			h.Discard()
			continue
		}

		for block := range slices.Chunk(h.MMap(), int(n)) {
			l := bytes.Count(block, []byte{'\n'})
			e := heap.Entropy(block)

			// add possibly remaining line
			if block[len(block)-1] != '\n' {
				l++
			}

			if e >= cmd.Min && e <= cmd.Max {
				title := text.Hide(h.String())
				start := text.Hide(fmt.Sprintf("@%d", off))

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
