package time

import (
	"errors"
	"strings"
	"time"

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

Filter flags:
  -M, --mode=[m|a|c|b]     Filter using timestamp type
  -N, --min=TIME           Minimum record time (RFC3339)
  -X, --max=TIME           Maximum record time (RFC3339)

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

	// filter flags
	Mode string    `short:"M" enum:"m,a,c,b" default:"a"`
	Min  time.Time `short:"N"`
	Max  time.Time `short:"X"`

	// paths
	Paths []string `arg:"" optional:""`
}

func (cmd *Time) Validate() error {
	if !cmd.Min.IsZero() && !cmd.Max.IsZero() && cmd.Min.After(cmd.Max) {
		return errors.New("invalid range")
	}

	return nil
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
			var ts time.Time

			switch strings.ToLower(cmd.Mode) {
			default:
			case "m":
				ts = e.Mtime
			case "a":
				ts = e.Atime
			case "c":
				ts = e.Ctime
			case "b":
				ts = e.Btime
			}

			if !cmd.Min.IsZero() && ts.Before(cmd.Min) {
				continue // to soon
			}

			if !cmd.Max.IsZero() && ts.After(cmd.Max) {
				continue // to late
			}

			fox.Writer.Match(formats.Auto(e, cmd.Json, cmd.Jsonl), fox.Regexp)
		}

		h.Free()
	}

	return nil
}
