package text

import (
	"fmt"
	"log"
	"strings"

	"github.com/alecthomas/kong"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
)

var Usage = strings.TrimSpace(`
Prints file text contents.

fox text [FLAGS ...] <PATHS ...>

Flags:
  -m, --min=NUMBER         minimum string length (default 3)
  -x, --max=NUMBER         maximal string length (default 256)
  -w, --wtf[=LEVEL]        show string classifications (w/ww/www)
  -F, --find=CLASS,...     show only strings with class(es)
  -1, --first              show only strings first class
  -P, --print              show only classification list

Examples:
  $ fox text -rw sample.exe
`)

type Text struct {
	Min   uint     `short:"m" default:"3"`
	Max   uint     `short:"x" default:"256"`
	Wtf   int      `short:"w" type:"counter"`
	Find  []string `short:"F" sep:","`
	First bool     `short:"1" and:"first,wtf"`
	Print bool     `short:"P"`
	Paths []string `arg:"" type:"path" optional:""`
}

func (cmd *Text) Validate() error {
	if cmd.Min > cmd.Max {
		log.Fatalln("invalid range")
	}

	if len(cmd.Find) > 0 {
		cmd.Wtf = 3
		cmd.First = false
		for i := range cmd.Find {
			cmd.Find[i] = strings.ToLower(cmd.Find[i])
		}
	}

	return nil
}

func (cmd *Text) AfterApply(app *kong.Kong, _ kong.Vars) error {
	if cmd.Print {
		for _, cls := range text.GetClasses(3) {
			fmt.Printf("%s\n", cls)
		}

		// exit early
		app.Exit(0)
	}

	return nil
}

func (cmd *Text) Run(cli *cli.Globals) error {
	if cli.Help || (len(cmd.Paths) == 0 && !cmd.Print) {
		fmt.Print(Usage)
		return nil
	}

	hs := cli.Load(cmd.Paths)
	defer cli.Discard()

	for _, h := range hs.Get() {
		if (hs.Len() > 1 || cli.Verbose > 0) && !cli.NoFile {
			_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide(text.Header(h.String())))
		}

		for s := range h.Strings(
			cmd.Min,
			cmd.Max,
			cmd.Wtf,
			cmd.Find,
			cmd.First,
		) {
			if !cli.NoLine && cmd.Wtf > 0 {
				_, _ = fmt.Fprintf(cli.Stdout, "%s  %s  %s\n", text.Hide(s.Off), s.Str, text.Hide(s.Cls))
			} else if !cli.NoLine {
				_, _ = fmt.Fprintf(cli.Stdout, "%s  %s\n", text.Hide(s.Off), s.Str)
			} else {
				_, _ = fmt.Fprintf(cli.Stdout, "%s\n", s.Str)
			}
		}
	}

	return nil
}
