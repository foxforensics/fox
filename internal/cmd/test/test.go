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

	"github.com/cuhsat/fox/v4/internal/pkg/data/apis/vt"
	"github.com/cuhsat/fox/v4/internal/pkg/hash"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

var Usage = strings.TrimSpace(`
Test suspicious files.

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

const ExitAlert = 3

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
	if len(cmd.Paths) == 0 {
		return text.Usage(Usage)
	}

	if len(cmd.Key) == 0 {
		log.Fatalln("VirusTotal API key required")
	}

	vt.Verbose = cli.Verbose

	cli.NoConvert = true // forced

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	for h := range ch {
		if !cmd.Domain && !cmd.Url && !cmd.Ip {
			if !cli.NoPretty {
				text.Title(h.String())
			}

			output(vt.TestFileHash(hash.MustSum(types.SHA256, h.Bytes()), cmd.Key))
		} else {
			for _, v := range strings.Split(string(h.Bytes()), "\n") {
				if len(strings.TrimSpace(v)) > 0 {
					if !cli.NoPretty {
						text.Title(v)
					}

					switch {
					case cmd.Domain:
						output(vt.TestDomain(v, cmd.Key))
					case cmd.Url:
						output(vt.TestUrl(base64.StdEncoding.EncodeToString([]byte(v)), cmd.Key))
					case cmd.Ip:
						output(vt.TestIp(v, cmd.Key))
					}
				}
			}
		}

		h.Discard()
	}

	return nil
}

func output(res *vt.Result, err error) {
	if err != nil {
		log.Println(err)
	}

	if res != nil {
		verdict := fmt.Sprintf("%s (%d/%d)", res.Label, res.Bad, res.All)

		if res.Alert {
			verdict = text.AsWarn(verdict)
		}

		text.Print(verdict)
	}
}
