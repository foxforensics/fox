package ad

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/alecthomas/kong"
	"go.foxforensics.dev/bootkey/bootkey"
	"go.foxforensics.dev/hashdump/extract"

	cli "go.foxforensics.dev/fox/v4/internal/cmd"

	"go.foxforensics.dev/fox/v4/internal/pkg/table"
	"go.foxforensics.dev/fox/v4/internal/pkg/text"
)

var Usage = strings.TrimSpace(`
fox ad [FLAGS...] NTDS SYSTEM

Flags:
  -j, --json               Show accounts as JSON objects
  -J, --jsonl              Show accounts as JSON lines

Secrets flags:
  -L, --lookup             Lookup hashes in the rainbow tables
  -H, --history            Extract also the users hash history
      --only-lm            Extract only the LM hashes (hashcat: 3000)
      --only-nt            Extract only the NT hashes (hashcat: 1000)

Example: Show NTLM secrets
  $ fox ad -LH NTDS.dit SYSTEM

Example: Show account infos
  $ fox ad -j NTDS.dit SYSTEM
`)

type Account struct {
	extract.Account
}

func (a *Account) ToJSON() string {
	b, _ := json.MarshalIndent(a, "", "  ")
	return string(b)
}

func (a *Account) ToJSONL() string {
	b, _ := json.Marshal(a)
	return string(b)
}

func (a *Account) ToNTLM(history bool) string {
	var sb strings.Builder

	// append actual hashes
	sb.WriteString(fmt.Sprintf("%s:%d:%s:%s:::",
		a.SAMAccountName,
		a.RID,
		a.format(a.LMHash, extract.DefaultLM),
		a.format(a.NTHash, extract.DefaultNT),
	))

	// append hash histories
	if history {
		for i := range a.NTHashHistory {
			sb.WriteString(fmt.Sprintf("\n%s_history%d:%d:%s:%s:::",
				a.SAMAccountName,
				i,
				a.RID,
				a.format(a.LMHashHistory[i], extract.DefaultLM),
				a.format(a.NTHashHistory[i], extract.DefaultNT),
			))
		}
	}

	return sb.String()
}

func (a *Account) OnlyLM() string {
	return a.format(a.LMHash, extract.DefaultLM)
}

func (a *Account) OnlyNT() string {
	return a.format(a.NTHash, extract.DefaultNT)
}

func (a *Account) format(sum string, def []byte) string {
	if pwd := table.Lookup(sum); len(pwd) > 0 {
		return text.AsBold(pwd)
	}

	if sum == fmt.Sprintf("%x", def) {
		return text.AsGray(sum)
	}

	return sum
}

type Ad struct {
	Json  bool `short:"j" xor:"json,jsonl"`
	Jsonl bool `short:"J" xor:"json,jsonl"`

	// secrets flags
	Lookup  bool `short:"L"`
	History bool `short:"H"`
	OnlyLm  bool `long:"only-lm" xor:"only-lm,only-nt"`
	OnlyNt  bool `long:"only-nt" xor:"only-lm,only-nt"`

	// paths
	Paths []string `arg:"" optional:""`
}

func (cmd *Ad) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	var err error

	if cmd.OnlyLm || cmd.OnlyNt {
		cmd.Json = false
		cmd.Jsonl = false
	}

	if cmd.Lookup {
		err = table.Build()
	}

	return err
}

func (cmd *Ad) Run(cli *cli.Globals) error {
	if len(cmd.Paths) < 2 {
		return text.Usage(Usage)
	}

	ch := cli.Load(cmd.Paths, true)
	defer cli.Discard()

	f1 := <-ch
	defer f1.Discard()

	f2 := <-ch
	defer f2.Discard()

	key, err := bootkey.ReadData(f2.Reader())

	if err != nil {
		return err
	}

	if cli.Verbose > 1 {
		log.Printf("BootKey %x\n", key)
	}

	pek, usr, err := extract.Extract(f1.Bytes(), key)

	if err != nil {
		return err
	}

	if cli.Verbose > 1 {
		for i, k := range pek {
			log.Printf("PEK #%d %x\n", i, k)
		}
	}

	if !cli.NoPretty {
		text.Title(f1.String())
	}

	for _, v := range usr {
		text.Match(cmd.format(&Account{v}), cli.Regexp)
	}

	if cli.Verbose > 1 {
		log.Printf("found %d account(s)\n", len(usr))
	}

	return nil
}

func (cmd *Ad) format(a *Account) string {
	switch {
	case cmd.Jsonl:
		return text.ColorizeAs(a.ToJSONL(), "json")
	case cmd.Json:
		return text.ColorizeAs(a.ToJSON(), "json")
	case cmd.OnlyLm:
		return a.OnlyLM()
	case cmd.OnlyNt:
		return a.OnlyNT()
	default:
		return a.ToNTLM(cmd.History)
	}
}
