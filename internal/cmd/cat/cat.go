package cat

import (
	"log/slog"
	"strings"

	"github.com/alecthomas/kong"
	"go.foxforensics.eu/fox/v4/internal/cmd"
	"go.foxforensics.eu/fox/v4/internal/pkg"
	"go.foxforensics.eu/fox/v4/internal/pkg/cat/buffer"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/fox/v4/internal/sys/loader"
	"go.foxforensics.eu/fox/v4/internal/sys/writer"
)

var Usage = strings.TrimSpace(`
Usage: fox cat [FLAGS...] <PATHS...>

Flags:
  -u, --uniq               Show only unique lines
  -t, --text               Force output as text
  -x, --hex                Force output as hex

Filter flags:
  -F, --find=REGEX         Filter using regular expression
  -C, --context=LINES      Lines surrounding a match
  -B, --before=LINES       Lines leading before a match
  -A, --after=LINES        Lines trailing after a match

Example: Show occurrences in event logs
  $ fox cat -FWinlogon ./**/*.evtx

Example: Show MBR in canonical hex
  $ fox cat -L512b image.dd

Report bugs at: foxforensics.eu/issues
`)

type Cat struct {
	Uniq bool `short:"u"`
	Text bool `short:"t" xor:"text,hex"`
	Hex  bool `short:"x" xor:"text,hex"`

	// filter flags
	Context uint `short:"C"`
	Before  uint `short:"B"`
	After   uint `short:"A"`

	// paths
	Paths []string `arg:"" optional:""`

	// internal
	uniq *pkg.Unique `kong:"-"`
}

func (cmd *Cat) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	if cmd.Uniq {
		cmd.uniq = pkg.NewUnique()
	}

	if cmd.Context > 0 {
		cmd.Before = cmd.Context
		cmd.After = cmd.Context
	}

	return nil
}

func (cmd *Cat) Run(fox *cmd.Globals) error {
	cmd.Paths = append(cmd.Paths, fox.Paths...)

	if len(cmd.Paths) == 0 {
		sys.Usage(Usage)
		return nil
	}

	// apply command specific params
	fox.Query.Before = cmd.Before
	fox.Query.After = cmd.After

	ch, err := fox.Init(cmd.Paths, cmd.Text || cmd.Hex)

	if err != nil {
		return err
	}

	for h := range ch {
		var hint string

		if !fox.NoPretty {
			fox.Writer.FileHeader(h.String())
		}

		if h.Stage >= loader.Convert {
			hint = "json" // default
		}

		if (h.IsText() && !cmd.Hex) || cmd.Text {
			cmd.renderText(fox, h.Bytes(), hint)
		} else {
			cmd.renderHex(fox, h.Bytes())
		}

		h.Discard()
	}

	return nil
}

func (cmd *Cat) renderText(fox *cmd.Globals, b []byte, hint string) {
	for l := range buffer.Text(&buffer.TextContext{
		Parent: fox.Context,
		Data:   b,
		Hint:   hint,
	}, fox).Lines {
		s := l.String

		if cmd.uniq != nil && !cmd.uniq.Is(s) {
			continue // not unique
		}

		if fox.Regexp != nil && l.Line != buffer.Sep {
			s = writer.MarkMatch(s, fox.Regexp)
		}

		if !fox.NoPretty {
			fox.Writer.Write("%s %s", writer.AsGray(l.Line), s)
		} else if l.Line == buffer.Sep {
			fox.Writer.Write(writer.AsGray(buffer.Sep))
		} else {
			fox.Writer.Write(s)
		}
	}
}

func (cmd *Cat) renderHex(fox *cmd.Globals, b []byte) {
	lastHex, wasCut := "", false

	for l := range buffer.Hex(&buffer.HexContext{
		Parent: fox.Context,
		Data:   b,
		Pretty: !fox.NoPretty,
	}, fox).Lines {
		if fox.Regexp != nil {
			if ok, err := fox.Regexp.MatchString(l.Values); !ok {
				if err != nil {
					slog.Error(err.Error())
				}
				continue // not matched afterward
			}
		}

		l.Values = writer.MarkMatch(l.Values, fox.Regexp)

		// cut similar lines for better readability
		if !fox.NoPretty && l.Values == lastHex {
			if !wasCut {
				wasCut = true
				fox.Writer.Write(writer.AsGray("*"))
			}
			continue
		}

		if !fox.NoPretty {
			fox.Writer.Write("%s  %s%s", writer.AsGray(l.Address), l.Values, l.String)
		} else {
			fox.Writer.Write(l.Values)
		}

		lastHex, wasCut = l.Values, false
	}
}
