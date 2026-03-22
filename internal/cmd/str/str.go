package str

import (
	"log"
	"strings"

	"github.com/alecthomas/kong"
	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/text/carver"
)

var Usage = strings.TrimSpace(`
Show file string contents.

fox str [FLAGS...] <PATHS...>

Flags:
  -n, --min=LENGTH         Minimum string length (default: 3)
  -x, --max=LENGTH         Maximal string length (default: 256)
  -a, --ascii              Show only strings with ASCII encoding
  -s, --sort               Sort strings alphabetically

Class flags:
  -w, --wtf[=LEVEL]        Show string classifications (w/ww/www)
  -F, --find=CLASS,...     Show only strings that match class(es)
  -1, --first              Show only strings first class
  -L, --list               Show only classification list

Examples:
  $ fox str -w sample.exe
`)

type Str struct {
	Min   uint `short:"n" default:"3"`
	Max   uint `short:"x" default:"256"`
	Ascii bool `short:"a"`
	Sort  bool `short:"s"`

	// class
	Wtf   int      `short:"w" type:"counter"`
	Find  []string `short:"F" sep:","`
	First bool     `short:"1" and:"first,wtf"`
	List  bool     `short:"L"`

	// paths
	Paths []string `arg:"" optional:""`
}

func (cmd *Str) Validate() error {
	if cmd.Min > cmd.Max {
		log.Fatalln("invalid range")
	}

	if (len(cmd.Find) > 0 || cmd.First) && cmd.Wtf == 0 {
		log.Fatalln("wtf required")
	}

	return nil
}

func (cmd *Str) AfterApply(app *kong.Kong, _ kong.Vars) error {
	if cmd.List {
		for _, cls := range carver.Classes(3) {
			text.Write(cls)
		}

		// exit early
		app.Exit(0)
	}

	return nil
}

func (cmd *Str) Run(cli *cli.Globals) error {
	if len(cmd.Paths)+len(cli.Paths) == 0 && !cmd.List {
		return text.Usage(Usage)
	}

	ch := cli.Load(cmd.Paths, true)
	defer cli.Discard()

	for h := range ch {
		if !cli.NoPretty {
			text.Title(h.String())
		}

		for l := range carver.New(&carver.Options{
			Min:      cmd.Min,
			Max:      cmd.Max,
			Ascii:    cmd.Ascii,
			Sort:     cmd.Sort,
			Wtf:      cmd.Wtf,
			Find:     cmd.Find,
			First:    cmd.First,
			Parallel: cli.Parallel,
		}).Carve(h.Bytes()) {
			if cli.Regexp != nil && !cli.Regexp.MatchString(l.Value) {
				continue // not matched afterward
			}

			l.Value = text.MarkMatch(l.Value, cli.Regexp)

			if !cli.NoPretty && len(l.Class) > 0 {
				text.Write("%s  %s [%s]", text.AsGray(l.Address), l.Value, text.AsBold(l.Class))
			} else if !cli.NoPretty {
				text.Write("%s  %s", text.AsGray(l.Address), l.Value)
			} else if len(l.Class) > 0 {
				text.Write("%s [%s]", l.Value, l.Class)
			} else {
				text.Write(l.Value)
			}
		}

		h.Discard()
	}

	return nil
}
