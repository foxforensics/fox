package info

import (
	"bytes"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/alecthomas/kong"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/data/api"
	"github.com/cuhsat/fox/v4/internal/pkg/data/api/vt"
	"github.com/cuhsat/fox/v4/internal/pkg/hash"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heap"
)

var Usage = strings.TrimSpace(`
Show file infos with verdict.

fox info [FLAGS...] <PATHS...>

Flags:
  -s, --sort               Sort files by path (slower)
  -b, --block=SIZE         Block size for analysis

Filter flags:
  -n, --min=VALUE          Minimum entropy value (default: 0.0)
  -x, --max=VALUE          Maximal entropy value (default: 1.0)

Format flags:
  -H, --human              Format size in human-readable units

Examples:
  $ fox info -n0.8 ./**/*

Remarks:
  Files hashes will be checked with VirusTotal, if FOX_API_KEY env is set.
`)

type Info struct {
	Sort  bool    `short:"s"`
	Block string  `short:"b"`
	Min   float64 `short:"n" default:"0.0"`
	Max   float64 `short:"x" default:"1.0"`

	// format
	Human bool `short:"H"`

	// paths
	Paths []string `arg:"" name:"path" optional:""`

	// hidden
	Key  string `hidden:"" long:"api-key"`
	Jack string `hidden:"" xor:"jack,john"`
	John string `hidden:"" xor:"jack,john"`

	// internal
	block uint64 `kong:"-"`
}

func (cmd *Info) Validate() error {
	if cmd.Min > cmd.Max {
		log.Fatalln("invalid range")
	}

	return nil
}

func (cmd *Info) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	switch {
	case len(cmd.Jack) > 0:
		cmd.Key = api.Decrypt(vt.ReserveKey1, cmd.Jack)

	case len(cmd.John) > 0:
		cmd.Key = api.Decrypt(vt.ReserveKey2, cmd.John)
	}

	if len(cmd.Block) > 0 {
		cmd.block = uint64(text.Mechanize(cmd.Block))
	}

	return nil
}

func (cmd *Info) Run(cli *cli.Globals) error {
	if len(cmd.Paths)+len(cli.Paths) == 0 {
		return text.Usage(Usage)
	}

	if cmd.Sort {
		cli.Parallel = 1 // single threaded
	}

	if !cli.NoPretty {
		text.Title(fmt.Sprintf("%-13s %11s %11s %s  %16s", "Entropy", "Lines", "Bytes", "Modified", "File"))
	}

	ch := cli.LoadPlain(cmd.Paths)
	defer cli.Discard()

	for h := range ch {
		var t = time.UnixMilli(int64(h.Time)).UTC().Format(time.RFC3339)
		var n = cmd.block
		var ver string
		var off int

		if n == 0 {
			n = h.Size
		}

		if h.Size == 0 {
			if cmd.Min == 0 {
				cmd.renderEmpty(h)
			}

			h.Discard()
			continue
		}

		if len(cmd.Key) > 0 {
			ver = cmd.verdict(h)
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
				start := fmt.Sprintf("[%08x]", off)

				if cmd.Human {
					size = text.Humanize(int64(len(block)))
				}

				if cmd.block > 0 {
					text.Write("%.10fe %10dl %11s %s  %s %s%s", e, l, size, t, start, h.String(), ver)
				} else {
					text.Write("%.10fe %10dl %11s %s  %s%s", e, l, size, t, h.String(), ver)
				}
			}

			off += len(block)
		}

		h.Discard()
	}

	return nil
}

func (cmd *Info) renderEmpty(h *heap.Heap) {
	t := time.UnixMilli(int64(h.Time)).UTC().Format(time.RFC3339)
	s := h.String()

	if cmd.block > 0 {
		s = "[00000000] " + s
	}

	text.Write(text.AsGray(fmt.Sprintf("%.10fe %10dl %10db %s  %s", 0.0, 0, 0, t, s)))
}

func (cmd *Info) renderBlock(h *heap.Heap, b []byte, o int) {
	t := time.UnixMilli(int64(h.Time)).UTC().Format(time.RFC3339)
	l := bytes.Count(b, []byte{'\n'})
	e := heap.Entropy(b)

	// add possibly remaining line
	if b[len(b)-1] != '\n' {
		l++
	}

	size := fmt.Sprintf("%db", len(b))
	start := fmt.Sprintf("[%08x]", o)

	if cmd.Human {
		size = text.Humanize(int64(len(b)))
	}

	if cmd.block > 0 {
		text.Write("%.10fe %10dl %11s %s  %s %s%s", e, l, size, t, start, h.String(), ver)
	} else {
		text.Write("%.10fe %10dl %11s %s  %s%s", e, l, size, t, h.String(), ver)
	}
}

func (cmd *Info) verdict(h *heap.Heap) string {
	res, err := vt.CheckFile(hash.MustSum(types.SHA256, h.Bytes()), cmd.Key)

	if err != nil {
		log.Println(err)
	}

	switch res.Verdict {
	case api.Unknown:
		return fmt.Sprintf(" [%s]", res.Verdict)
	case api.Unrated, api.Clean:
		return fmt.Sprintf(" [%s]", text.AsGray(res.Verdict))
	default:
		return fmt.Sprintf(" [%s]", text.AsWarn(res.Verdict))
	}
}
