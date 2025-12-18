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

	hs := cli.Bootstrap(cmd.Paths)
	defer cli.ThrowAway()

	for _, h := range hs.Get() {
		if hs.Len() > 1 && !cli.NoFile {
			_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide(text.Header(h.String())))
		}

		for l := range buffer.Text(h, 2).Lines {
			s := l.Str

			if cli.Filter != nil && l.Nr != buffer.Sep {
				s = text.MarkMatch(s, cli.Filter)
			}

			if !cli.NoLine && l.Nr == buffer.Sep {
				_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide(buffer.Sep))
			} else if !cli.NoLine {
				_, _ = fmt.Fprintf(cli.Stdout, "%s %s\n", text.Hide(l.Nr), s)
			} else {
				_, _ = fmt.Fprintf(cli.Stdout, "%s\n", s)
			}
		}
	}

	return nil
}
