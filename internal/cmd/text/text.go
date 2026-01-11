package text

import (
	"fmt"
	"log"
	"strings"

	"github.com/alecthomas/kong"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/text/carver"
)

var Usage = strings.TrimSpace(`
Prints file text contents.

fox text [FLAGS...] <PATHS...>

Flags:
  -m, --min=NUMBER         minimum string length (default: 6)
  -x, --max=NUMBER         maximal string length (default: 256)
  -s, --sort               sort strings alphabetically (slower)

Class:
  -w, --wtf[=LEVEL]        show string classifications (w/ww/www)
  -F, --find=CLASS,...     show only strings with classes
  -1, --first              show only strings first class
  -l, --list               show only classification list

Example:
  $ fox text -w sample.exe
`)

type Text struct {
	Min  uint `short:"m" default:"6"`
	Max  uint `short:"x" default:"256"`
	Sort bool `short:"s"`

	// class
	Wtf   int      `short:"w" type:"counter"`
	Find  []string `short:"F" sep:","`
	First bool     `short:"1" and:"first,wtf"`
	List  bool     `short:"l"`

	// paths
	Paths []string `arg:"" type:"path" optional:""`
}

func (cmd *Text) Validate() error {
	if cmd.Min > cmd.Max {
		log.Fatalln("invalid range")
	}

	return nil
}

func (cmd *Text) AfterApply(app *kong.Kong, _ kong.Vars) error {
	if cmd.List {
		for _, cls := range carver.Classes(3) {
			fmt.Printf("%s\n", cls)
		}

		// exit early
		app.Exit(0)
	}

	return nil
}

func (cmd *Text) Run(cli *cli.Globals) error {
	if cli.Help || (len(cmd.Paths) == 0 && !cmd.List) {
		fmt.Print(Usage)
		return nil
	}

	cli.NoConvert = true // forced

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	for h := range ch {
		if !cli.NoFile {
			_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide(text.Header(h.String())))
		}

		for s := range carver.New(&carver.Options{
			Min:     cmd.Min,
			Max:     cmd.Max,
			Sort:    cmd.Sort,
			Wtf:     cmd.Wtf,
			Find:    cmd.Find,
			First:   cmd.First,
			Profile: cli.Profile,
		}).Carve(h.MMap()) {
			if !cli.NoLine && cmd.Wtf > 0 {
				_, _ = fmt.Fprintf(cli.Stdout, "%s  %s  %s\n", text.Hide(s.Adr), s.Str, text.Hide(s.Cls))
			} else if !cli.NoLine {
				_, _ = fmt.Fprintf(cli.Stdout, "%s  %s\n", text.Hide(s.Adr), s.Str)
			} else {
				_, _ = fmt.Fprintf(cli.Stdout, "%s\n", s.Str)
			}
		}

		h.Discard()
	}

	return nil
}
