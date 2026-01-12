package test

import (
	"fmt"
	"log"
	"strings"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/data/extern/vt"
	"github.com/cuhsat/fox/v4/internal/pkg/hash"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

var Usage = strings.TrimSpace(`
Prints file test results.

fox test [FLAGS...] <PATHS...>

Flags:
  -l, --level=[=LEVEL]     VirusTotal report level (l/ll/lll)
  -k, --key=KEY            VirusTotal API key

Example:
  $ fox test -l sample.exe
`)

type Test struct {
	Level int      `short:"l" long:"level" type:"counter"`
	Key   string   `short:"k" long:"key"`
	Paths []string `arg:"" name:"path" type:"path" optional:""`
}

func (cmd *Test) Validate() error {
	if len(cmd.Key) == 0 {
		log.Fatalln("key required")
	}

	return nil
}

func (cmd *Test) Run(cli *cli.Globals) error {
	if cli.Help || len(cmd.Paths) == 0 {
		fmt.Print(Usage)
		return nil
	}

	cli.NoConvert = true // forced

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	for h := range ch {
		var res string
		var err error

		if !cli.NoFile {
			_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide(text.Header(h.String())))
		}

		sum := hash.MustSum(types.SHA256, h.MMap())

		switch {
		case cmd.Level > 1:
			res, err = vt.GetReport(sum, cmd.Key, cli.NoPretty)
		case cmd.Level > 0:
			res, err = vt.GetResult(sum, cmd.Key)
		default:
			res, err = vt.GetLabel(sum, cmd.Key)
		}

		if err != nil {
			if cli.Filter != nil && !cli.Filter.MatchString(res) {
				continue // not matched afterward
			}

			res = text.MarkMatch(res, cli.Filter)

			_, _ = fmt.Fprintln(cli.Stdout, res)
		} else {
			log.Println(err)
		}

		h.Discard()
	}

	return nil
}
