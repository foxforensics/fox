package test

import (
	"fmt"
	"log"
	"strings"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/data/extern/virus"
	"github.com/cuhsat/fox/v4/internal/pkg/hash"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

var Usage = strings.TrimSpace(`
Prints file test results.

fox test [FLAGS...] <PATHS...>

Flags:
  -k, --key=APIKEY         Set key for VirusTotal API

Example:
  $ fox test sample.exe
`)

type Test struct {
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

	if cli.Verbose > 2 {
		virus.Trace = true
	}

	cli.NoConvert = true // forced

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	for h := range ch {
		if !cli.NoFile {
			_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide(text.Header(h.String())))
		}

		res, err := virus.Test(hash.MustSum(types.SHA256, h.MMap()), cmd.Key)

		if err != nil {
			log.Println(err)
		}

		for _, r := range res {
			if r.Alert {
				_, _ = fmt.Fprintf(cli.Stdout, "%s %s\n", text.Warn(r.Result), text.Hide(r.Engine))
			} else {
				_, _ = fmt.Fprintf(cli.Stdout, "%s %s\n", r.Result, text.Hide(r.Engine))
			}
		}

		h.Discard()
	}

	return nil
}
