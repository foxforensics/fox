package str

import (
	"log"
	"strings"

	"github.com/alecthomas/kong"

	cli "go.foxforensics.dev/fox/v4/internal/cmd"

	"go.foxforensics.dev/fox/v4/internal/pkg/text"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/carver"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/lookup"
)

var Usage = strings.TrimSpace(`
fox str [FLAGS...] <PATHS...>

Flags:
  -n, --min=LENGTH         Minimum string length (default: 3)
  -x, --max=LENGTH         Maximal string length (default: 256)
  -a, --ascii              Show only strings with ASCII encoding
  -s, --sort               Sort strings alphabetically
  -m, --trim               Trim strings whitespaces

Class flags:
  -w, --what[=LEVEL]       Show string classifications (w/ww/www)
  -F, --find=CLASS,...     Show only strings that match class(es)
  -1, --first              Show only strings first class
      --list               Show only classification list

Lookup flags:
  -L, --lookup             Lookup URLs, IPs and domains via VirusTotal

Remarks:
  A VirusTotal API key is required for lookup. 

Example: Show only long ASCII strings
  $ fox str -ant8 sample.exe

Example: Show all URLs in a binary
  $ fox str -wFurl sample.exe
`)

type Str struct {
	Min   uint `short:"n" default:"3"`
	Max   uint `short:"x" default:"256"`
	Ascii bool `short:"a"`
	Sort  bool `short:"s"`
	Trim  bool `short:"m"`

	// class flags
	What  int      `short:"w" type:"counter"`
	Find  []string `short:"F" sep:","`
	First bool     `short:"1" and:"first,what"`
	List  bool

	// lookup flags
	Lookup bool `short:"L"`

	// paths
	Paths []string `arg:"" optional:""`
}

func (cmd *Str) Validate() error {
	if cmd.Min > cmd.Max {
		log.Fatalln("invalid range")
	}

	if (len(cmd.Find) > 0 || cmd.First) && cmd.What == 0 {
		log.Fatalln("what required")
	}

	return nil
}

func (cmd *Str) AfterApply(app *kong.Kong, _ kong.Vars) error {
	if cmd.List {
		db := text.BuildDB(3)

		for _, s := range db.List() {
			text.Write(s)
		}

		// exit early
		app.Exit(0)
	}

	if cmd.Lookup && cmd.What == 0 {
		cmd.What = 1
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

		for str := range carver.New(&carver.Options{
			Min:      cmd.Min,
			Max:      cmd.Max,
			Ascii:    cmd.Ascii,
			Sort:     cmd.Sort,
			Trim:     cmd.Trim,
			What:     cmd.What,
			Find:     cmd.Find,
			First:    cmd.First,
			Parallel: cli.Parallel,
		}).Carve(h.Bytes()) {
			if cli.Regexp != nil && !cli.Regexp.MatchString(str.Value) {
				continue // not matched afterward
			}

			str.Value = text.MarkMatch(str.Value, cli.Regexp)

			if !cli.NoPretty && cmd.Lookup && lookup.Lookup(str, cli.Verbose) {
				text.Write("%s  %s [%s]", text.AsGray(str.Address), text.AsWarn(str.Value), text.AsBold(str.Classes))
			} else if !cli.NoPretty && len(str.Classes) > 0 {
				text.Write("%s  %s [%s]", text.AsGray(str.Address), str.Value, text.AsBold(str.Classes))
			} else if !cli.NoPretty {
				text.Write("%s  %s", text.AsGray(str.Address), str.Value)
			} else if len(str.Classes) > 0 {
				text.Write("%s [%s]", str.Value, str.Classes)
			} else {
				text.Write(str.Value)
			}
		}

		h.Discard()
	}

	return nil
}
