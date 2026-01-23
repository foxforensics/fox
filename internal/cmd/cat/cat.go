package cat

import (
	"fmt"
	"strings"

	"github.com/alecthomas/kong"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types/buffer"
)

var Usage = strings.TrimSpace(`
Prints contents.

fox cat [FLAGS...] <PATHS...>

Flags:
  -C, --context=NUMBER     lines surrounding context of a match
  -B, --before=NUMBER      lines leading context before a match
  -A, --after=NUMBER       lines trailing context after a match

Example:
  $ fox -eWinlogon ./**/*.evtx
`)

type Cat struct {
	Context uint     `short:"C"`
	Before  uint     `short:"B"`
	After   uint     `short:"A"`
	Paths   []string `arg:"" type:"path" optional:""`
}

func (cmd *Cat) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	if cmd.Context > 0 {
		cmd.Before = cmd.Context
		cmd.After = cmd.Context
	}

	return nil
}

func (cmd *Cat) Run(cli *cli.Globals) error {
	if len(cmd.Paths)+len(cli.File) == 0 {
		fmt.Print(Usage)
		return nil
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

		for l := range buffer.Text(h, cli).Lines {
			s := l.String

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
