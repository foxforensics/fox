package time

import (
	"strings"

	"go.foxforensics.eu/fox/v5/internal/cmd"
	"go.foxforensics.eu/fox/v5/internal/cmd/time/parser"
	"go.foxforensics.eu/fox/v5/internal/pkg"
	"go.foxforensics.eu/fox/v5/library/formats"
)

var Usage = strings.TrimSpace(`
Usage: fox time [FLAGS...] <PATHS...>

Flags:
  -s, --sort               Sort timeline chronologically
  -j, --json               Show timeline as JSON objects
  -l, --jsonl              Show timeline as JSON lines

Example: Show MFT entries as bodyfile
  $ fox time ./$MFT

Example: Show entries chronologically
  $ fox time -s ./**/*.pf

Report bugs at: foxforensics.eu/issues
`)

type Time struct {
	Sort  bool `short:"s"`
	Json  bool `short:"j" xor:"json,jsonl"`
	Jsonl bool `short:"l" xor:"json,jsonl"`

	// paths
	Paths []string `arg:"" optional:""`
}

func (cmd *Time) Run(fox *cmd.Globals) error {
	cmd.Paths = append(cmd.Paths, fox.Paths...)

	if len(cmd.Paths) == 0 {
		pkg.Usage(Usage)
		return nil
	}

	heaps, err := fox.Init(cmd.Paths, true)

	if err != nil {
		return err
	}

	for h := range heaps {
		for e := range parser.New(&parser.Options{
			Sort:    cmd.Sort,
			Threads: fox.Threads,
		}).Parse(fox.Context, h.Bytes()) {
			fox.Writer.Match(formats.Auto(e, cmd.Json, cmd.Jsonl), fox.Regexp)
		}

		h.Free()
	}

	return nil
}
