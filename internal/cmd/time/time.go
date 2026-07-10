package time

import (
	"fmt"
	"log/slog"
	"strings"

	"go.foxforensics.eu/fox/v4/internal/cmd"
	"go.foxforensics.eu/fox/v4/internal/pkg/time/entry"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/fox/v4/internal/sys/writer"
	"go.foxforensics.eu/fox/v4/library/binaries/bin/lnk"
	"go.foxforensics.eu/fox/v4/library/binaries/bin/mft"
	"go.foxforensics.eu/fox/v4/library/binaries/bin/pf"
	"go.foxforensics.eu/fox/v4/library/formats"
)

var Usage = strings.TrimSpace(`
Usage: fox time [FLAGS...] <PATHS...>

Flags:
  -s, --sort               Sort timeline chronologically
  -j, --json               Show timeline as JSON objects
  -J, --jsonl              Show timeline as JSON lines

Format flags:
  -b, --bodyfile           Show in Body file format
  -t, --timesketch         Show in Timesketch format

Example: Show MFT entries as body file
  $ fox time -b ./$MFT

Report bugs at: foxforensics.eu/issues
`)

type Time struct {
	Sort  bool `short:"s"`
	Json  bool `short:"j" xor:"json,jsonl"`
	Jsonl bool `short:"J" xor:"json,jsonl"`

	// format flags
	Bodyfile   bool `short:"b" xor:"bodyfile,timesketch"`
	Timesketch bool `short:"t" xor:"bodyfile,timesketch"`

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
		switch {
		case mft.Detect(h.Bytes()):
			slog.Debug("file detected as mft")

			for _, e := range mft.Parse(h.Bytes()) {
				fox.Writer.Match(cmd.format(e), fox.Regexp)
			}

		case lnk.Detect(h.Bytes()):
			slog.Debug("file detected as lnk")

			for _, e := range lnk.Parse(h.Bytes()) {
				fox.Writer.Match(cmd.format(e), fox.Regexp)
			}

		case pf.Detect(h.Bytes()):
			slog.Debug("file detected as pf")

			for _, e := range pf.Parse(h.Bytes()) {
				fox.Writer.Match(cmd.format(e), fox.Regexp)
			}

		default:
			slog.Debug(fmt.Sprintf("file not supported %s", h))
		}

		h.Free()
	}

	//slices.SortStableFunc(src, func(a, b *event.Event) int {
	//	return strings.Compare(a.SortKey(), b.SortKey())
	//})

	return nil
}

func (cmd *Time) format(e *entry.Entry) string {
	switch {
	case cmd.Json:
		return writer.ColorizeAs(formats.AsJSON(e), "json")
	case cmd.Jsonl:
		return writer.ColorizeAs(formats.AsJSONL(e), "json")
	case cmd.Bodyfile:
		return e.AsBodyfile()
	case cmd.Timesketch:
		return e.AsTimesketch()
	case e.Anomaly:
		return writer.AsBold(e.String())
	default:
		return e.String()
	}
}
