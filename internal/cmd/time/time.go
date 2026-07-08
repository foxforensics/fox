package time

import (
	"fmt"
	"log/slog"
	"strings"

	"go.foxforensics.eu/fox/v4/internal/cmd"
	"go.foxforensics.eu/fox/v4/internal/lib/binaries/bin/mft"
	"go.foxforensics.eu/fox/v4/internal/lib/formats"
	"go.foxforensics.eu/fox/v4/internal/pkg/time/entry"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/fox/v4/internal/sys/writer"
)

var Usage = strings.TrimSpace(`
Usage: fox time [FLAGS...] <PATHS...>

Flags:
  -c, --csv                Show timeline as CSV list
  -j, --json               Show timeline as JSON objects
  -J, --jsonl              Show timeline as JSON lines

Example: Show MFT entries as body file
  $ fox time "$MFT"s

Report bugs at: foxforensics.eu/issues
`)

type Time struct {
	Csv   bool `short:"c" xor:"csv,json,jsonl"`
	Json  bool `short:"j" xor:"csv,json,jsonl"`
	Jsonl bool `short:"J" xor:"csv,json,jsonl"`

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

			for _, e := range mft.Parse(h.Bytes()) {
				fox.Writer.Match(cmd.format(&e), fox.Regexp)
			}
		}

		h.Free()
	}

	return nil
}

func (cmd *Time) format(e *entry.Entry) string {
	switch {
	case cmd.Csv:
		return e.AsCSV()

	case cmd.Json, cmd.Jsonl:
		return formats.Auto(e, cmd.Json, cmd.Jsonl)

	case e.Anomaly:
		return writer.AsBold(e.String())

	default:
		return e.String()
	}
}
