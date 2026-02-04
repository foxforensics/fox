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
Tests suspicious files.

fox test [FLAGS...] PATHS...

Flags:
  -D, --domain             File(s) contains domains
  -U, --url                File(s) contains urls
  -I, --ip                 File(s) contains ips

Required:
  -k, --key=APIKEY         VirusTotal API key

Examples:
  $ fox test ioc.exe
`)

type Test struct {
	Key    string   `short:"k"`
	Domain bool     `short:"D" xor:"domain,url,ip"`
	Url    bool     `short:"U" xor:"domain,url,ip"`
	Ip     bool     `short:"I" xor:"domain,url,ip"`
	Paths  []string `arg:"" type:"path" optional:""`
}

func (cmd *Test) Run(cli *cli.Globals) error {
	var alert bool

	if len(cmd.Paths) == 0 {
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

	for h := range ch {
		if !(cmd.Domain || cmd.Url || cmd.Ip) {
			res, err := vt.TestFileHash(hash.MustSum(types.SHA256, h.Bytes()), cmd.Key)
			alert = alert || cmd.output(cli, res, err, h.String())
		} else {
			var res *vt.Result
			var err error

			for _, v := range strings.Split(string(h.Bytes()), "\n") {
				if len(v) == 0 {
					continue
				}

				switch {
				case cmd.Domain:
					res, err = vt.TestDomain(v, cmd.Key)
				case cmd.Url:
					res, err = vt.TestUrl(base64.StdEncoding.EncodeToString([]byte(v)), cmd.Key)
				case cmd.Ip:
					res, err = vt.TestIp(v, cmd.Key)
				}

				if err != nil {
					log.Println(err)
				}

				if res != nil {
					alert = alert || cmd.output(cli, res, err, v)
				}
			}
		}

		h.Discard()
	}

	if alert {
		cli.Exit(3)
	}

	return nil
}

func (cmd *Test) output(cli *cli.Globals, res *vt.Result, err error, h string) bool {
	if !cli.NoFile {
		_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide(text.Header(h)))
	}

	if err != nil {
		log.Println(err)
		return false
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

	return res.Alert
}
