package str

import (
	"errors"
	"log"
	"strings"

	"github.com/alecthomas/kong"

	cli "go.foxforensics.eu/fox/v4/internal/cmd"

	"go.foxforensics.eu/fox/v4/internal/pkg/text"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/carver"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/lookup"
)

var Usage = strings.TrimSpace(`
Usage: fox str [FLAGS...] <list|PATHS...>

Flags:
  -l, --lookup             Lookup strings via VirusTotal (IP/DNS/URL) 
  -a, --ascii              Show only strings with ASCII encoding
  -s, --sort               Sort strings alphabetically
  -t, --trim               Trim strings whitespaces
  -N, --min=LENGTH         Minimum string length (default: 3)
  -X, --max=LENGTH         Maximal string length (default: 256)

Class flags:
  -w, --what[=LEVEL]       Show string classifications (w/ww/www)
  -C, --class=NAME,...     Show only classes that match name(es)

Remarks:
  If 'list' is specified as path, only the built-in classifications
  will be shown. A VirusTotal API key is required for lookup.

Example: Show only long ASCII strings
  $ fox str -atN8 sample.exe

Example: Show all URLs in a binary
  $ fox str -wCurl sample.exe

Report bugs at: foxforensics.eu/issues
`)

type Str struct {
	Lookup bool `short:"l"`
	Ascii  bool `short:"a"`
	Sort   bool `short:"s"`
	Trim   bool `short:"t"`
	Min    uint `short:"N" default:"3"`
	Max    uint `short:"X" default:"256"`

	// class flags
	What  int      `short:"w" type:"counter"`
	Class []string `short:"C" sep:","`

	// paths
	Paths []string `arg:"" optional:""`
}

func (cmd *Str) Validate() error {
	if cmd.Min > cmd.Max {
		return errors.New("invalid range")
	}

	if cmd.Lookup {
		log.Println("warning: data will be transmitted to a third-party service!")
	}

	return nil
}

func (cmd *Str) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	if cmd.Lookup && cmd.What == 0 {
		cmd.What = 1
	}

	if len(cmd.Class) > 0 {
		cmd.What = 3
	}

	return nil
}

func (cmd *Str) Run(cli *cli.Globals) error {
	cmd.Paths = append(cmd.Paths, cli.Input...)

	if len(cmd.Paths) == 0 {
		return text.Usage(Usage)
	}

	if cmd.Paths[0] == "list" {
		db := text.BuildDB(3)

		for _, s := range db.List() {
			text.Stdout.Write(s)
		}

		// exit early
		return nil
	}

	ch := cli.Load(cmd.Paths, true)
	defer cli.Discard()

	for h := range ch {
		if !cli.NoPretty {
			text.Stdout.Header(h.String())
		}

		for str := range carver.New(&carver.Options{
			Min:     cmd.Min,
			Max:     cmd.Max,
			Ascii:   cmd.Ascii,
			Sort:    cmd.Sort,
			Trim:    cmd.Trim,
			What:    cmd.What,
			Class:   cmd.Class,
			Threads: cli.Threads,
		}).Carve(h.Bytes()) {
			if cli.Regexp != nil {
				if ok, _ := cli.Regexp.MatchString(str.Value); !ok {
					continue // not matched afterward
				}
			}

			str.Value = text.MarkMatch(str.Value, cli.Regexp)

			if !cli.NoPretty && cmd.Lookup && lookup.Lookup(str, cli.Verbose) {
				text.Stdout.Write("%s  %s [%s]", text.AsGray(str.Address), text.AsWarn(str.Value), text.AsBold(str.Classes))
			} else if !cli.NoPretty && len(str.Classes) > 0 {
				text.Stdout.Write("%s  %s [%s]", text.AsGray(str.Address), str.Value, text.AsBold(str.Classes))
			} else if !cli.NoPretty {
				text.Stdout.Write("%s  %s", text.AsGray(str.Address), str.Value)
			} else if len(str.Classes) > 0 {
				text.Stdout.Write("%s [%s]", str.Value, str.Classes)
			} else {
				text.Stdout.Write(str.Value)
			}
		}

		h.Discard()
	}

	return nil
}
