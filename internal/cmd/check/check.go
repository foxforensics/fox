package check

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/xxtea/xxtea-go/xxtea"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/data/api"
	"github.com/cuhsat/fox/v4/internal/pkg/data/api/hibp"
	"github.com/cuhsat/fox/v4/internal/pkg/data/api/vt"
	"github.com/cuhsat/fox/v4/internal/pkg/hash"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

var Usage = strings.TrimSpace(`
Check files, domains, mails, URLs and IPs.

fox check [FLAGS...] PATHS...

Flags:
  -j, --json               Show results as JSON objects
  -J, --jsonl              Show results as JSON lines

Content flags:
  -D, --domain             File(s) contains domains
  -M, --mail               File(s) contains mails
  -U, --url                File(s) contains urls
  -I, --ip                 File(s) contains ips

Required:
      --hp-key=APIKEY      API key for HaveIBeenPwned
      --vt-key=APIKEY      API key for VirusTotal

Examples:
  $ fox check sample.exe
`)

// Encrypted backup keys for emergency use
const (
	VtKey1 = "47ba3c085f105fff4fa186ce769f8a35f98bc3010fd8e25c9a90c1bf70696120b9fe1a5c6328bf0deae4eebdcc9f5df156a27efd923eaad648f3e8ab26fcc8f6753233b8"
	VtKey2 = "44201ef4cbffe7edd1a7d2279a1fc3019700c3620da45d0542014b8a7be0fd7b53125c3e474c6db7360f4f538d56bfe15bd416b0d2a77c02a37d0ffc5015694b41c9f117"
)

type Check struct {
	Json  bool `short:"j" xor:"json,jsonl"`
	Jsonl bool `short:"J" xor:"json,jsonl"`

	// content flags
	Domain bool `short:"D" xor:"domain,main,url,ip"`
	Mail   bool `short:"M" xor:"domain,main,url,ip"`
	Url    bool `short:"U" xor:"domain,main,url,ip"`
	Ip     bool `short:"I" xor:"domain,main,url,ip"`

	// required
	HpKey string `long:"hp-key"`
	VtKey string `long:"vt-key"`

	// hidden
	Reserve1 string `short:"1" xor:"reserve1,reserve2" hidden:""`
	Reserve2 string `short:"2" xor:"reserve1,reserve2" hidden:""`

	// paths
	Paths []string `arg:"" optional:""`
}

func (cmd *Check) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	switch {
	case len(cmd.Reserve1) > 0:
		v, _ := hex.DecodeString(VtKey1)
		cmd.VtKey = string(xxtea.Decrypt(v, []byte(cmd.Reserve1)))

	case len(cmd.Reserve2) > 0:
		v, _ := hex.DecodeString(VtKey2)
		cmd.VtKey = string(xxtea.Decrypt(v, []byte(cmd.Reserve2)))
	}

	return nil
}

func (cmd *Check) Run(cli *cli.Globals) error {
	if len(cmd.Paths) == 0 {
		return text.Usage(Usage)
	}

	if len(cmd.HpKey) == 0 && cmd.Mail {
		log.Fatalln("HaveIBeenPwned API key required")
	}

	if len(cmd.VtKey) == 0 && !cmd.Mail {
		log.Fatalln("VirusTotal API key required")
	}

	cli.NoConvert = true // forced

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	for h := range ch {
		if !cmd.Domain && !cmd.Mail && !cmd.Url && !cmd.Ip {
			res, err := vt.CheckFileHash(hash.MustSum(types.SHA256, h.Bytes()), cmd.VtKey)

			if err != nil {
				log.Println(err)
			}

			if res != nil {
				text.Print(cmd.format(res, h.String()))
			}
		} else {
			for _, v := range strings.Split(string(h.Bytes()), "\n") {
				if len(strings.TrimSpace(v)) > 0 {
					var res *api.Result
					var err error

					switch {
					case cmd.Domain:
						res, err = vt.CheckDomain(v, cmd.VtKey)
					case cmd.Mail:
						res, err = hibp.CheckMail(v, cmd.HpKey)
					case cmd.Url:
						res, err = vt.CheckUrl(base64.StdEncoding.EncodeToString([]byte(v)), cmd.VtKey)
					case cmd.Ip:
						res, err = vt.CheckIp(v, cmd.VtKey)
					}

					if err != nil {
						log.Println(err)
						continue
					}

					if res != nil {
						text.Print(cmd.format(res, v))
					}
				}
			}
		}

		h.Discard()
	}

	return nil
}

func (cmd *Check) format(res *api.Result, src string) string {
	var line string

	switch {
	case cmd.Jsonl:
		line = text.ColorizeStringAs(res.ToJSONL(), "json")
	case cmd.Json:
		line = text.ColorizeStringAs(res.ToJSON(), "json")
	case res.Verdict == api.Unknown:
		line = fmt.Sprintf("%s:%s", src, res.Verdict)
	case res.Verdict == api.Unrated:
		line = fmt.Sprintf("%s:%s", src, text.AsGray(res.Verdict))
	case res.Verdict == api.Clean:
		line = fmt.Sprintf("%s:%s", src, text.AsGray(res.Verdict))
	case res.Verdict == api.Suspicious:
		line = fmt.Sprintf("%s:%s", src, text.AsWarn(res.Verdict))
	default:
		line = fmt.Sprintf("%s:%s", src, text.AsWarn(res.Verdict))
	}

	return line
}
