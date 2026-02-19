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
  -n, --min=LENGTH         minimum string length (default: 3)
  -x, --max=LENGTH         maximal string length (default: 256)
  -a, --ascii              shows only strings with ASCII encoding
  -s, --sort               sorts strings alphabetically

Class flags:
  -w, --wtf[=LEVEL]        shows string classifications (w/ww/www)
  -F, --find=CLASS,...     shows only strings that match class(es)
  -1, --first              shows only strings first class
  -L, --list               shows only classification list

Format flags:
  -D, --decimal            format addresses as decimals

Examples:
  $ fox text -w ioc.exe
`)

type Text struct {
	Min   uint `short:"n" default:"3"`
	Max   uint `short:"x" default:"256"`
	Ascii bool `short:"a"`
	Sort  bool `short:"s"`

	// class
	Wtf   int      `short:"w" type:"counter"`
	Find  []string `short:"F" sep:","`
	First bool     `short:"1" and:"first,wtf"`
	List  bool     `short:"L"`

	// format
	Decimal bool `short:"D"`

	// paths
	Paths []string `arg:"" optional:""`
}

func (cmd *Text) Validate() error {
	if cmd.Min > cmd.Max {
		log.Fatalln("invalid range")
	}

	if (len(cmd.Find) > 0 || cmd.First) && cmd.Wtf == 0 {
		log.Fatalln("wtf required")
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
	if len(cmd.Paths)+len(cli.Paths) == 0 && !cmd.List {
		fmt.Println(Usage)
		return nil
	}

	cli.NoConvert = true // forced

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	for h := range ch {
		if !cli.NoFile {
			_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide(text.Title(h.String())))
		}

		for l := range carver.New(&carver.Options{
			Min:      cmd.Min,
			Max:      cmd.Max,
			Ascii:    cmd.Ascii,
			Sort:     cmd.Sort,
			Wtf:      cmd.Wtf,
			Find:     cmd.Find,
			First:    cmd.First,
			Decimal:  cmd.Decimal,
			Parallel: cli.Threads,
		}).Carve(h.Bytes()) {
			if cli.Regexp != nil && !cli.Regexp.MatchString(l.Value) {
				continue // not matched afterward
			}

			l.Value = text.MarkMatch(l.Value, cli.Regexp)

			if !cli.NoLine && cmd.Wtf > 0 {
				_, _ = fmt.Fprintf(cli.Stdout, "%s  %s  %s\n", text.Hide(l.Address), l.Value, text.Hide(l.Class))
			} else if !cli.NoLine {
				_, _ = fmt.Fprintf(cli.Stdout, "%s  %s\n", text.Hide(l.Address), l.Value)
			} else {
				_, _ = fmt.Fprintf(cli.Stdout, "%s\n", l.Value)
			}
		}

		h.Discard()
	}

	return nil
}
