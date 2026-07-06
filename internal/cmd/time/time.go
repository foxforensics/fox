package time

import (
	"fmt"
	"log/slog"
	"strings"

	"go.foxforensics.eu/fox/v4/internal/cmd"
	"go.foxforensics.eu/fox/v4/internal/lib/binaries/bin/mft"
	"go.foxforensics.eu/fox/v4/internal/lib/formats"
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

	for h := range ch {
		if mft.Detect(h.Bytes()) {
			slog.Debug(fmt.Sprintf("mft detected %s", h))

			for _, v := range mft.ToBody(h.Bytes()) {
				fox.Writer.Match(formats.Auto(v, cmd.Json, cmd.Jsonl), fox.Regexp)
			}
		}

		h.Free()
	}

	return nil
}
