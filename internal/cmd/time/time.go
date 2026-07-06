package time

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"go.foxforensics.eu/fox/v4/internal/cmd"
	"go.foxforensics.eu/fox/v4/internal/lib/binaries/bin/mft"
	"go.foxforensics.eu/fox/v4/internal/lib/formats"
	"go.foxforensics.eu/fox/v4/internal/pkg/time/body"
	"go.foxforensics.eu/fox/v4/internal/sys"
)

var Usage = strings.TrimSpace(`
Usage: fox time [FLAGS...] <PATHS...>

Flags:
  -b, --body               Show timeline as body file
  -j, --json               Show timeline as JSON objects
  -J, --jsonl              Show timeline as JSON lines

Example: Show MFT timeline
  $ fox time "$MFT"

Report bugs at: foxforensics.eu/issues
`)

type Time struct {
	Body  bool `short:"b" xor:"body,json,jsonl"`
	Json  bool `short:"j" xor:"body,json,jsonl"`
	Jsonl bool `short:"J" xor:"body,json,jsonl"`

	// paths
	Paths []string `arg:"" optional:""`
}

func (cmd *Time) Run(fox *cmd.Globals) error {
	cmd.Paths = append(cmd.Paths, fox.Paths...)

	if len(cmd.Paths) == 0 {
		sys.Usage(Usage)
		return nil
	}

	ch, err := fox.Init(cmd.Paths, true)

	if err != nil {
		return err
	}

	bf := make([]body.Body, 0)

	for h := range ch {
		if mft.Detect(h.Bytes()) {
			slog.Debug(fmt.Sprintf("mft detected %s", h))

			bf = append(bf, mft.ToBody(h.Bytes())...)
		}

		h.Free()
	}

	slices.SortStableFunc(bf, func(a, b body.Body) int {
		return strings.Compare(a.Name, b.Name)
	})

	for _, b := range bf {
		fox.Writer.Match(formats.Auto(b, cmd.Json, cmd.Jsonl), fox.Regexp)
	}

	return nil
}
