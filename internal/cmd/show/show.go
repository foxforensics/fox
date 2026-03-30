package show

import (
	"github.com/alecthomas/kong"

	cli "go.foxforensics.dev/fox/v4/internal/cmd"

	"go.foxforensics.dev/fox/v4/internal/pkg/text"
	"go.foxforensics.dev/fox/v4/internal/pkg/text/unique"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/buffer"
)

type Show struct {
	// unique
	Uniq bool    `short:"u" xor:"uniq,dist"`
	Dist float64 `short:"D" xor:"uniq,dist"`

	// filter
	Context uint `short:"C"`
	Before  uint `short:"B"`
	After   uint `short:"A"`

	// display
	ForceText bool `short:"T" xor:"force-text,force-hex" long:"force-text"`
	ForceHex  bool `short:"X" xor:"force-text,force-hex" long:"force-hex"`

	// paths
	Paths []string `arg:"" optional:""`

	// internal
	uniq unique.Unique `kong:"-"`
}

func (cmd *Show) AfterApply(_ *kong.Kong, _ kong.Vars) error {
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

func (cmd *Show) Run(cli *cli.Globals) error {
	ch := cli.Load(cmd.Paths, false)

	if cmd.ForceText || cmd.ForceHex {
		cli.NoConvert = true
	}

	// apply command specific context
	cli.Filter.Before = cmd.Before
	cli.Filter.After = cmd.After

	defer cli.Discard()

	for h := range ch {
		if !cli.NoPretty {
			text.Title(h.String())
		}

		if (h.IsText() && !cmd.ForceHex) || cmd.ForceText {
			cmd.renderText(cli, h.Bytes(), h.Hint)
		} else {
			cmd.renderHex(cli, h.Bytes())
		}

		h.Discard()
	}

	return nil
}

func (cmd *Show) renderText(cli *cli.Globals, b []byte, hint string) {
	for l := range buffer.Text(cli, &buffer.TextContext{
		Data: b,
		Hint: hint,
	}).Lines {
		s := l.String

		if cmd.uniq != nil && !cmd.uniq.IsUnique(s) {
			continue // not unique
		}

		if cli.Regexp != nil && l.Line != buffer.Sep {
			s = text.MarkMatch(s, cli.Regexp)
		}

		if !cli.NoPretty {
			text.Write("%s %s", text.AsGray(l.Line), s)
		} else if l.Line == buffer.Sep {
			text.Write(text.AsGray(buffer.Sep))
		} else {
			text.Write(s)
		}
	}
}

func (cmd *Show) renderHex(cli *cli.Globals, b []byte) {
	lastHex, wasCut := "", false

	for l := range buffer.Hex(cli, &buffer.HexContext{
		Data:   b,
		Pretty: !cli.NoPretty,
	}).Lines {
		if cli.Regexp != nil && !cli.Regexp.MatchString(l.Values) {
			continue // not matched afterward
		}

		l.Values = text.MarkMatch(l.Values, cli.Regexp)

		// cut similar lines for better readability
		if !cli.NoPretty && l.Values == lastHex {
			if !wasCut {
				wasCut = true
				text.Write(text.AsGray("*"))
			}
			continue
		}

		if !cli.NoPretty {
			text.Write("%s  %s%s", text.AsGray(l.Address), text.MarkZero(l.Values), l.String)
		} else {
			text.Write(l.Values)
		}

		lastHex, wasCut = l.Values, false
	}
}
