package test

import (
	"encoding/base64"
	"fmt"
	"log"
	"math"
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
  -k, --key=APIKEY         Set key for VirusTotal API
  -U, --url=URL            Test suspicious URL
  -I, --ip=IP              Test suspicious IP

Example:
  $ fox test sample.exe
`)

type Test struct {
	Key   string   `short:"k"`
	Url   string   `short:"U"`
	Ip    string   `short:"I"`
	Paths []string `arg:"" name:"path" type:"path" optional:""`
}

func (cmd *Test) Validate() error {
	if len(cmd.Key) == 0 {
		log.Fatalln("key required")
	}

	return nil
}

func (cmd *Test) Run(cli *cli.Globals) error {
	if cli.Help || len(cmd.Paths)+len(cmd.Ip)+len(cmd.Url) == 0 {
		fmt.Print(Usage)
		return nil
	}

	if cli.Verbose > 2 {
		vt.Trace = true
	}

	cli.NoConvert = true // forced

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	if len(cmd.Ip) > 0 {
		res, err := vt.TestIp(cmd.Ip, cmd.Key)
		cmd.output(cli, res, err, cmd.Ip)
	}

	if len(cmd.Url) > 0 {
		res, err := vt.TestUrl(base64.StdEncoding.EncodeToString([]byte(cmd.Url)), cmd.Key)
		cmd.output(cli, res, err, cmd.Url)
	}

	for h := range ch {
		res, err := vt.TestHash(hash.MustSum(types.SHA256, h.MMap()), cmd.Key)
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

	n := int(math.Log10(float64(len(res)))) + 1

	for i, r := range res {
		nr := fmt.Sprintf("%0*d", n, i+1)

		if r.Alert {
			r.Result = text.Warn(r.Result)
		}

		if !cli.NoLine {
			_, _ = fmt.Fprintf(cli.Stdout, "%s %s %s\n", text.Hide(nr), r.Result, text.Hide(r.Engine))
		} else {
			_, _ = fmt.Fprintf(cli.Stdout, "%s %s\n", r.Result, text.Hide(r.Engine))
		}
	}
}
