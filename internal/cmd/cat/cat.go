package cat

import (
	"errors"
	"fmt"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types/buffer"
)

type Cat struct {
	Paths []string `arg:"" type:"path" optional:""`
}

func (cmd *Cat) Run(cli *cli.Globals) error {
	if len(cmd.Paths) == 0 {
		return errors.New("path required")
	}

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	for h := range ch {
		if !cli.NoFile {
			_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide(text.Header(h.String())))
		}

		for l := range buffer.Text(h, cli.Profile).Lines {
			s := l.Str

			if cli.Filter != nil && l.Nr != buffer.Sep {
				s = text.MarkMatch(s, cli.Filter)
			}

			if !cli.NoLine && l.Nr == buffer.Sep {
				_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide(text.Line()))
			} else if l.Nr == buffer.Sep {
				_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide("--"))
			} else if !cli.NoLine {
				_, _ = fmt.Fprintf(cli.Stdout, "%s %s\n", text.Hide(l.Nr), s)
			} else {
				_, _ = fmt.Fprintf(cli.Stdout, "%s\n", s)
			}
		}

		h.Discard()
	}

	return nil
}
