package cat

import (
	"strings"

	"github.com/alecthomas/kong"
	"go.foxforensics.eu/fox/v4/internal/cmd"
	"go.foxforensics.eu/fox/v4/internal/pkg/types"
	buffer2 "go.foxforensics.eu/fox/v4/internal/pkg/types/buffer"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/fox/v4/internal/sys/terminal"
)

var Usage = strings.TrimSpace(`
Usage: fox cat [FLAGS...] <PATHS...>

Flags:
  -t, --text               Force output as text
  -x, --hex                Force output as hex

Unique flags:
  -u, --uniq               Unique by XXH3 hash sum
  -D, --dist=LENGTH        Unique by Levenshtein distance

Filter flags:
  -F, --find=PATTERN       Filter using regular expression
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
	Text bool `short:"t" xor:"text,hex"`
	Hex  bool `short:"x" xor:"text,hex"`

	// filter flags
	Uniq    bool    `short:"u" xor:"uniq,dist"`
	Dist    float64 `short:"D" xor:"uniq,dist"`
	Context uint    `short:"C"`
	Before  uint    `short:"B"`
	After   uint    `short:"A"`

	// paths
	Paths []string `arg:"" optional:""`

	// internal
	uniq *types.Unique `kong:"-"`
}

func (cmd *Cat) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	switch {
	case cmd.Uniq:
		cmd.uniq = types.NewUnique(types.Hash)
	case cmd.Dist > 0:
		cmd.uniq = types.NewUnique(types.Distance)
		cmd.uniq.SetLimit(cmd.Dist)
	}

	if cmd.Context > 0 {
		cmd.Before = cmd.Context
		cmd.After = cmd.Context
	}

	return nil
}

func (cmd *Cat) Run(fox *cmd.Globals) error {
	cmd.Paths = append(cmd.Paths, fox.Input...)

	if len(cmd.Paths) == 0 {
		return sys.Usage(Usage)
	}

	ch, err := fox.Init(cmd.Paths, cmd.Text || cmd.Hex)

	if err != nil {
		return err
	}

	defer fox.Discard()

	// apply command specific context
	fox.Filters.Before = cmd.Before
	fox.Filters.After = cmd.After

	for h := range ch {
		if !fox.NoPretty {
			sys.Stdout.Header(h.String())
		}

		if (h.IsText() && !cmd.Hex) || cmd.Text {
			cmd.renderText(fox, h.Bytes(), h.Hint)
		} else {
			cmd.renderHex(fox, h.Bytes())
		}

		h.Discard()
	}

	return nil
}

func (cmd *Cat) renderText(fox *cmd.Globals, b []byte, hint string) {
	for l := range buffer2.Text(&buffer2.TextContext{
		Parent: fox.Context,
		Data:   b,
		Hint:   hint,
	}, fox).Lines {
		s := l.String

		if cmd.uniq != nil && !cmd.uniq.IsUnique(s) {
			continue // not unique
		}

		if fox.Regexp != nil && l.Line != buffer2.Sep {
			s = terminal.MarkMatch(s, fox.Regexp)
		}

		if !fox.NoPretty {
			sys.Stdout.Write("%s %s", terminal.AsGray(l.Line), s)
		} else if l.Line == buffer2.Sep {
			sys.Stdout.Write(terminal.AsGray(buffer2.Sep))
		} else {
			sys.Stdout.Write(s)
		}
	}
}

func (cmd *Cat) renderHex(fox *cmd.Globals, b []byte) {
	lastHex, wasCut := "", false

	for l := range buffer2.Hex(&buffer2.HexContext{
		Parent: fox.Context,
		Data:   b,
		Pretty: !fox.NoPretty,
	}, fox).Lines {
		if fox.Regexp != nil {
			if ok, _ := fox.Regexp.MatchString(l.Values); !ok {
				continue // not matched afterward
			}
		}

		l.Values = terminal.MarkMatch(l.Values, fox.Regexp)

		// cut similar lines for better readability
		if !fox.NoPretty && l.Values == lastHex {
			if !wasCut {
				wasCut = true
				sys.Stdout.Write(terminal.AsGray("*"))
			}
			continue
		}

		if !fox.NoPretty {
			sys.Stdout.Write("%s  %s%s", terminal.AsGray(l.Address), l.Values, l.String)
		} else {
			sys.Stdout.Write(l.Values)
		}

		lastHex, wasCut = l.Values, false
	}
}
