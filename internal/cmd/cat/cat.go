package cat

import (
	"fmt"
	"strings"

	"github.com/alecthomas/kong"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/text/unique"
	"github.com/cuhsat/fox/v4/internal/pkg/types/buffer"
)

var Usage = strings.TrimSpace(`
Prints file contents.

fox cat [FLAGS...] <PATHS...>

Filter Flags:
  -u, --uniq               filters output by unique hash (XXH3)
  -D, --dist=LENGTH        filters output by Levenshtein distance (slow)
  -e, --regexp=PATTERN     filters output by regular expression

RegExp flags:
  -C, --context=LINES      lines surrounding context of a match
  -B, --before=LINES       lines leading context before a match
  -A, --after=LINES        lines trailing context after a match

Examples:
  $ fox -eWinlogon ./**/*.evtx
`)

type Cat struct {
	// filter
	Uniq bool    `short:"u" xor:"uniq,dist"`
	Dist float64 `short:"D" xor:"uniq,dist"`

	// regexp
	Context uint `short:"C"`
	Before  uint `short:"B"`
	After   uint `short:"A"`

	// paths
	Paths []string `arg:"" type:"path" optional:""`

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
	if len(cmd.Paths)+len(cli.Paths) == 0 {
		fmt.Print(Usage)
		return nil
	}

	if cmd.Dist > 0 {
		cli.NoColor = true
	}

	ch := cli.Load(cmd.Paths)

	// apply command specific context
	cli.Filter.Before = cmd.Before
	cli.Filter.After = cmd.After

	defer cli.Discard()

	for h := range ch {
		if !cli.NoFile {
			_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide(text.Header(h.String())))
		}

		for l := range buffer.Text(h, cli, new(buffer.TextContext)).Lines {
			s := l.String

			if cmd.uniq != nil && !cmd.uniq.IsUnique(s) {
				continue // not unique
			}

			if cli.Regexp != nil && l.Line != buffer.Sep {
				s = text.MarkMatch(s, cli.Regexp)
			}

			if !cli.NoLine && l.Line == buffer.Sep {
				_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide(text.Line()))
			} else if l.Line == buffer.Sep {
				_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide("--"))
			} else if !cli.NoLine {
				_, _ = fmt.Fprintf(cli.Stdout, "%s %s\n", text.Hide(l.Line), s)
			} else {
				_, _ = fmt.Fprintf(cli.Stdout, "%s\n", s)
			}
		}

		h.Discard()
	}

	return nil
}
