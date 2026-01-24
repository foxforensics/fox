package test

import (
	"encoding/base64"
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
Prints test results.

fox test [FLAGS...] [PATHS...]

Flags:
  -k, --key=APIKEY         Sets VirusTotal API key
  -U, --url=URL,...        Tests suspicious URL
  -I, --ip=IP,...          Tests suspicious IP

Examples:
  $ fox test sample.exe
`)

type Test struct {
	Key   string   `short:"k"`
	Url   []string `short:"U" sep:","`
	Ip    []string `short:"I" sep:","`
	Paths []string `arg:"" name:"path" type:"path" optional:""`
}

func (cmd *Test) Run(cli *cli.Globals) error {
	if len(cmd.Paths)+len(cli.Paths)+len(cmd.Ip)+len(cmd.Url) == 0 {
		fmt.Print(Usage)
		return nil
	}

	if len(cmd.Key) == 0 {
		log.Fatalln("VirusTotal API key required")
	}

	if cli.Verbose > 2 {
		vt.Trace = true
	}

	cli.NoConvert = true // forced

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	for _, v := range cmd.Ip {
		res, err := vt.TestIp(v, cmd.Key)
		cmd.output(cli, res, err, v)
	}

	for _, v := range cmd.Url {
		res, err := vt.TestUrl(base64.StdEncoding.EncodeToString([]byte(v)), cmd.Key)
		cmd.output(cli, res, err, v)
	}

	for h := range ch {
		res, err := vt.TestHash(hash.MustSum(types.SHA256, h.Bytes()), cmd.Key)
		cmd.output(cli, res, err, h.String())
		h.Discard()
	}

	return nil
}

func (cmd *Test) output(cli *cli.Globals, res []vt.Entry, err error, h string) {
	if !cli.NoFile {
		_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide(text.Header(h)))
	}

	if err != nil {
		log.Println(err)
	}

	for _, r := range res {
		if r.Alert {
			r.Result = text.Warn(r.Result)
		}

		_, _ = fmt.Fprintf(cli.Stdout, "%s  %s\n", r.Result, text.Hide(r.Engine))
	}
}
