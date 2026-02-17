package test

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/xxtea/xxtea-go/xxtea"

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
  -D, --domain             file(s) contains domains
  -U, --url                file(s) contains urls
  -I, --ip                 file(s) contains ips

Required:
  -k, --key=APIKEY         VirusTotal API key

Examples:
  $ fox test ioc.exe
`)

// Encrypted backup keys for emergency use
const (
	Key1 = "47ba3c085f105fff4fa186ce769f8a35f98bc3010fd8e25c9a90c1bf70696120b9fe1a5c6328bf0deae4eebdcc9f5df156a27efd923eaad648f3e8ab26fcc8f6753233b8"
	Key2 = "44201ef4cbffe7edd1a7d2279a1fc3019700c3620da45d0542014b8a7be0fd7b53125c3e474c6db7360f4f538d56bfe15bd416b0d2a77c02a37d0ffc5015694b41c9f117"
)

type Test struct {
	Domain bool `short:"D" xor:"domain,url,ip"`
	Url    bool `short:"U" xor:"domain,url,ip"`
	Ip     bool `short:"I" xor:"domain,url,ip"`

	// required
	Key string `short:"k"`

	// hidden
	One string `xor:"one,two" hidden:""`
	Two string `xor:"one,two" hidden:""`

	// paths
	Paths []string `arg:"" optional:""`
}

func (cmd *Test) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	switch {
	case len(cmd.One) > 0:
		v, _ := hex.DecodeString(Key1)
		cmd.Key = string(xxtea.Decrypt(v, []byte(cmd.One)))

	case len(cmd.Two) > 0:
		v, _ := hex.DecodeString(Key2)
		cmd.Key = string(xxtea.Decrypt(v, []byte(cmd.Two)))
	}

	return nil
}

func (cmd *Test) Run(cli *cli.Globals) error {
	var alert bool

	if len(cmd.Paths) == 0 {
		fmt.Println(Usage)
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
		_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide(text.Title(h)))
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
