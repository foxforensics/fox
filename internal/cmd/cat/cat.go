package cat

import (
	"strings"

	"github.com/alecthomas/kong"
	buffer2 "go.foxforensics.eu/fox/v4/internal/cmd/cat/buffer"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/unique"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/fox/v4/internal/sys/output"

	cli "go.foxforensics.eu/fox/v4/internal/cmd"
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
	uniq unique.Unique `kong:"-"`
}

func (cmd *Cat) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	switch {
	case cmd.Uniq:
		cmd.uniq = unique.ByHash()
	case cmd.Dist > 0:
		cmd.uniq = unique.ByDistance(cmd.Dist)
	}

	if cmd.Context > 0 {
		cmd.Before = cmd.Context
		cmd.After = cmd.Context
	}

	return nil
}

func (cmd *Cat) Run(cli *cli.Globals) error {
	cmd.Paths = append(cmd.Paths, cli.Input...)

	if len(cmd.Paths) == 0 {
		return sys.Usage(Usage)
	}

	ch, err := cli.Init(cmd.Paths, false)

	if err != nil {
		return err
	}

	defer cli.Discard()

	if cmd.Text || cmd.Hex {
		cli.NoConvert = true
	}

	// apply command specific context
	cli.Filters.Before = cmd.Before
	cli.Filters.After = cmd.After

	for h := range ch {
		if !cli.NoPretty {
			sys.Stdout.Header(h.String())
		}

		if (h.IsText() && !cmd.Hex) || cmd.Text {
			cmd.renderText(cli, h.Bytes(), h.Hint)
		} else {
			cmd.renderHex(cli, h.Bytes())
		}

		h.Discard()
	}

	return nil
}

func (cmd *Cat) renderText(cli *cli.Globals, b []byte, hint string) {
	for l := range buffer2.Text(&buffer2.TextContext{
		Parent: cli.Context,
		Data:   b,
		Hint:   hint,
	}, cli).Lines {
		s := l.String

		if cmd.uniq != nil && !cmd.uniq.IsUnique(s) {
			continue // not unique
		}

		if cli.Regexp != nil && l.Line != buffer2.Sep {
			s = output.MarkMatch(s, cli.Regexp)
		}

		if !cli.NoPretty {
			sys.Stdout.Write("%s %s", output.AsGray(l.Line), s)
		} else if l.Line == buffer2.Sep {
			sys.Stdout.Write(output.AsGray(buffer2.Sep))
		} else {
			sys.Stdout.Write(s)
		}
	}
}

func (cmd *Cat) renderHex(cli *cli.Globals, b []byte) {
	lastHex, wasCut := "", false

	for l := range buffer2.Hex(&buffer2.HexContext{
		Parent: cli.Context,
		Data:   b,
		Pretty: !cli.NoPretty,
	}, cli).Lines {
		if cli.Regexp != nil {
			if ok, _ := cli.Regexp.MatchString(l.Values); !ok {
				continue // not matched afterward
			}
		}

		l.Values = output.MarkMatch(l.Values, cli.Regexp)

		// cut similar lines for better readability
		if !cli.NoPretty && l.Values == lastHex {
			if !wasCut {
				wasCut = true
				sys.Stdout.Write(output.AsGray("*"))
			}
			continue
		}

		if !cli.NoPretty {
			sys.Stdout.Write("%s  %s%s", output.AsGray(l.Address), l.Values, l.String)
		} else {
			sys.Stdout.Write(l.Values)
		}

		lastHex, wasCut = l.Values, false
	}
}
