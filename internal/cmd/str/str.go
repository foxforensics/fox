package str

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
	"go.foxforensics.eu/fox/v4/internal/cmd"
	"go.foxforensics.eu/fox/v4/internal/pkg/str"
	"go.foxforensics.eu/fox/v4/internal/pkg/str/carver"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/fox/v4/internal/sys/writer"
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
	cmd.Paths = append(cmd.Paths, fox.Paths...)

	if len(cmd.Paths) == 0 {
		sys.Usage(Usage)
		return nil
	}

	if cmd.Paths[0] == "list" {
		db := str.BuildDB(3)

		for _, s := range db.List() {
			fmt.Println(s)
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
			fox.Writer.Header(h.String())
		}

		for str := range carver.New(&carver.Options{
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

			str.Value = writer.MarkMatch(str.Value, fox.Regexp)

			if !fox.NoPretty && len(str.Classes) > 0 {
				fox.Writer.Write("%s  %s [%s]", writer.AsGray(str.Address), str.Value, writer.AsBold(str.Classes))
			} else if !fox.NoPretty {
				fox.Writer.Write("%s  %s", writer.AsGray(str.Address), str.Value)
			} else if len(str.Classes) > 0 {
				fox.Writer.Write("%s [%s]", str.Value, str.Classes)
			} else {
				fox.Writer.Write(str.Value)
			}
		}

		h.Discard()
	}

	return nil
}
