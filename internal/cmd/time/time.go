package time

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"go.foxforensics.eu/fox/v4/internal/cmd"
	"go.foxforensics.eu/fox/v4/internal/pkg/time/entry"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/fox/v4/internal/sys/heap"
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

Example: Show entries chronologically
  $ fox time -s ./**/*.pf

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

	heaps, err := fox.Init(cmd.Paths, true)

	if err != nil {
		return err
	}

	ch := cmd.parse(fox.Context, heaps)

	if cmd.Sort {
		ch = cmd.sort(fox.Context, ch)
	}

	for e := range ch {
		fox.Writer.Match(cmd.format(e), fox.Regexp)
	}

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

func (cmd *Time) parse(ctx context.Context, heaps <-chan *heap.Heap) <-chan *entry.Entry {
	entries := make(chan *entry.Entry, 4096)

	go func() {
		defer close(entries)

		var parse func(b []byte) []entry.Entry

		for h := range heaps {
			switch {
			case mft.Detect(h.Bytes()):
				slog.Debug("file detected as mft")
				parse = mft.Parse

			case lnk.Detect(h.Bytes()):
				slog.Debug("file detected as lnk")
				parse = lnk.Parse

			case pf.Detect(h.Bytes()):
				slog.Debug("file detected as pf")
				parse = pf.Parse

			default:
				slog.Debug(fmt.Sprintf("file not supported %s", h))
				h.Free()
				continue
			}

			for _, e := range parse(h.Bytes()) {
				select {
				case entries <- &e:
				case <-ctx.Done():
					h.Free()
					return
				}
			}

			h.Free()
		}
	}()

	return entries
}

func (cmd *Time) sort(ctx context.Context, ch <-chan *entry.Entry) <-chan *entry.Entry {
	sorted := make(chan *entry.Entry, cap(ch))

	go func() {
		defer close(sorted)

		v := make([]*entry.Entry, 0)

		for e := range ch {
			v = append(v, e)
		}

		slices.SortStableFunc(v, func(a, b *entry.Entry) int {
			return a.SortKey().Compare(b.SortKey())
		})

		for _, e := range v {
			select {
			case sorted <- e:
			case <-ctx.Done():
				return
			}
		}
	}()

	return sorted
}
