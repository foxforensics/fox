package str

import (
	"errors"
	"strings"

	"github.com/alecthomas/kong"
	"go.foxforensics.eu/fox/v4/internal/cmd"
	carver2 "go.foxforensics.eu/fox/v4/internal/pkg/types/carver"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/fox/v4/internal/sys/terminal"
)

var Usage = strings.TrimSpace(`
Usage: fox str [FLAGS...] <list|PATHS...>

Flags:
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
  will be shown.

Example: Show only long ASCII strings
  $ fox str -atN8 sample.exe

Example: Show all URLs in a binary
  $ fox str -wCurl sample.exe

Report bugs at: foxforensics.eu/issues
`)

type Str struct {
	Ascii bool `short:"a"`
	Sort  bool `short:"s"`
	Trim  bool `short:"t"`
	Min   uint `short:"N" default:"3"`
	Max   uint `short:"X" default:"256"`

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

	return nil
}

func (cmd *Str) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	if len(cmd.Class) > 0 {
		cmd.What = 3
	}

	return nil
}

func (cmd *Str) Run(fox *cmd.Globals) error {
	cmd.Paths = append(cmd.Paths, fox.Input...)

	if len(cmd.Paths) == 0 {
		return sys.Usage(Usage)
	}

	if cmd.Paths[0] == "list" {
		db := carver2.BuildDB(3)

		for _, s := range db.List() {
			fox.Stdout.Write(s)
		}

		// exit early
		return nil
	}

	ch, err := fox.Init(cmd.Paths, true)

	if err != nil {
		return err
	}

	for h := range ch {
		if !fox.NoPretty {
			fox.Stdout.Header(h.String())
		}

		for str := range carver2.New(&carver2.Options{
			Min:     cmd.Min,
			Max:     cmd.Max,
			Ascii:   cmd.Ascii,
			Sort:    cmd.Sort,
			Trim:    cmd.Trim,
			What:    cmd.What,
			Class:   cmd.Class,
			Threads: fox.Threads,
		}).Carve(fox.Context, h.Bytes()) {
			if fox.Regexp != nil {
				if ok, _ := fox.Regexp.MatchString(str.Value); !ok {
					continue // not matched afterward
				}
			}

			str.Value = terminal.MarkMatch(str.Value, fox.Regexp)

			if !fox.NoPretty && len(str.Classes) > 0 {
				fox.Stdout.Write("%s  %s [%s]", terminal.AsGray(str.Address), str.Value, terminal.AsBold(str.Classes))
			} else if !fox.NoPretty {
				fox.Stdout.Write("%s  %s", terminal.AsGray(str.Address), str.Value)
			} else if len(str.Classes) > 0 {
				fox.Stdout.Write("%s [%s]", str.Value, str.Classes)
			} else {
				fox.Stdout.Write(str.Value)
			}
		}

		h.Discard()
	}

	return nil
}
