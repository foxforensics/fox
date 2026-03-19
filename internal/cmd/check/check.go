package check

import (
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/xxtea/xxtea-go/xxtea"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/data/api"
	"github.com/cuhsat/fox/v4/internal/pkg/data/api/vt"
	"github.com/cuhsat/fox/v4/internal/pkg/hash"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

var Usage = strings.TrimSpace(`
Check suspicious items online.

fox check [FLAGS...] PATHS...

Flags:
  -j, --json               Show results as JSON objects
  -J, --jsonl              Show results as JSON lines

Required:
      --key=APIKEY         API key for VirusTotal

Examples:
  $ fox check sample.exe
`)

// Encrypted backup keys for emergency use
const (
	Key1 = "47ba3c085f105fff4fa186ce769f8a35f98bc3010fd8e25c9a90c1bf70696120b9fe1a5c6328bf0deae4eebdcc9f5df156a27efd923eaad648f3e8ab26fcc8f6753233b8"
	Key2 = "44201ef4cbffe7edd1a7d2279a1fc3019700c3620da45d0542014b8a7be0fd7b53125c3e474c6db7360f4f538d56bfe15bd416b0d2a77c02a37d0ffc5015694b41c9f117"
)

type Check struct {
	Json  bool `short:"j" xor:"json,jsonl"`
	Jsonl bool `short:"J" xor:"json,jsonl"`

	// required
	ApiKey string `long:"api-key"`

	// hidden
	Reserve1 string `short:"1" xor:"reserve1,reserve2" hidden:""`
	Reserve2 string `short:"2" xor:"reserve1,reserve2" hidden:""`

	// paths
	Paths []string `arg:"" optional:""`
}

func (cmd *Check) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	switch {
	case len(cmd.Reserve1) > 0:
		v, _ := hex.DecodeString(Key1)
		cmd.ApiKey = string(xxtea.Decrypt(v, []byte(cmd.Reserve1)))

	case len(cmd.Reserve2) > 0:
		v, _ := hex.DecodeString(Key2)
		cmd.ApiKey = string(xxtea.Decrypt(v, []byte(cmd.Reserve2)))
	}

	return nil
}

func (cmd *Check) Run(cli *cli.Globals) error {
	if len(cmd.Paths) == 0 {
		return text.Usage(Usage)
	}

	if len(cmd.ApiKey) == 0 {
		log.Fatalln("VirusTotal API key required")
	}

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	for h := range ch {
		res, err := vt.CheckFile(hash.MustSum(types.SHA256, h.Bytes()), cmd.ApiKey)

		if err != nil {
			log.Println(err)
		}

		if res != nil {
			text.Write(cmd.format(res, h.String()))
		}

		h.Discard()
	}

	return nil
}

func (cmd *Check) format(res *api.Result, src string) string {
	var line string

	switch {
	case cmd.Jsonl:
		line = text.ColorizeAs(res.ToJSONL(), "json")
	case cmd.Json:
		line = text.ColorizeAs(res.ToJSON(), "json")
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
