package cat

import (
	"errors"
	"fmt"
	"log"

	"github.com/pbnjay/memory"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types/buffer"
)

const limit = 0.95

type Cat struct {
	Paths []string `arg:"" type:"path" optional:""`
}

func (cmd *Cat) Run(cli *cli.Globals) error {
	if len(cmd.Paths) == 0 {
		return errors.New("path required")
	}

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	if !cli.NoWarnings && float32(memory.FreeMemory()/memory.TotalMemory()) > limit {
		log.Println("warning: low memory may cause swapping!")
	}

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
