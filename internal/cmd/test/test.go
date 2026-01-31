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
Prints file test results.

fox test [FLAGS...] [PATHS...]

Flags:
  -D, --domain=NAME,...    Tests suspicious domain(s)
  -U, --url=URL,...        Tests suspicious url(s)
  -I, --ip=IP,...          Tests suspicious ip(s)

Required:
  -k, --key=APIKEY         Sets VirusTotal API key

Examples:
  $ fox test ioc.exe
`)

type Test struct {
	Key    string   `short:"k"`
	Domain []string `short:"D" sep:","`
	Url    []string `short:"U" sep:","`
	Ip     []string `short:"I" sep:","`
	Paths  []string `arg:"" name:"path" type:"path" optional:""`
}

func (cmd *Test) Run(cli *cli.Globals) error {
	if len(cmd.Paths)+len(cli.Paths)+len(cmd.Ip)+len(cmd.Url)+len(cmd.Domain) == 0 {
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

	for _, v := range cmd.Domain {
		res, err := vt.TestDomain(v, cmd.Key)
		cmd.output(cli, res, err, v)
	}

	for h := range ch {
		res, err := vt.TestFileHash(hash.MustSum(types.SHA256, h.Bytes()), cmd.Key)
		cmd.output(cli, res, err, h.String())
		h.Discard()
	}

	return nil
}

func (cmd *Test) output(cli *cli.Globals, res *vt.Result, err error, h string) {
	if !cli.NoFile {
		_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide(text.Header(h)))
	}

	if err != nil {
		log.Println(err)
		return
	}

	_, _ = fmt.Fprint(cli.Stdout, "VirusTotal:\n\n")

	for _, e := range res.Entries {
		if e.Alert {
			e.Result = text.Warn(e.Result)
		}

		e.Engine += strings.Repeat(text.Hide("."), 30-len(e.Engine))

		_, _ = fmt.Fprintf(cli.Stdout, "  %s %s\n", e.Engine, e.Result)
	}

	if len(res.Entries) > 0 {
		_, _ = fmt.Fprintln(cli.Stdout, "")
	}

	if res.Alert {
		res.Label = text.Warn(res.Label)
	}

	_, _ = fmt.Fprintf(cli.Stdout, "  (%d of %d) %s\n\n", res.Bad, res.All, text.Bold(res.Label))
}
