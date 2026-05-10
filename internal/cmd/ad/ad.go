package ad

import (
	"fmt"
	"log"
	"strings"

	"github.com/alecthomas/kong"
	"go.foxforensics.dev/bootkey/bootkey"
	"go.foxforensics.dev/hashdump/extract"

	cli "go.foxforensics.dev/fox/v4/internal/cmd"

	"go.foxforensics.dev/fox/v4/internal/pkg/text"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/record"
)

var Usage = strings.TrimSpace(`
fox ad [FLAGS...] NTDS SYSTEM

Account flags:
  -F, --find=hash          Show only accounts that match hash
  -j, --json               Show accounts as JSON objects
  -J, --jsonl              Show accounts as JSON lines

Secrets flags:
  -H, --history            Extract also the LM and NT hash history
      --lm                 Extract just the LM hashes (hashcat: 3000)
      --nt                 Extract just the NT hashes (hashcat: 1000)

Example: Show NTLM secrets
  $ fox ad -H NTDS.dit SYSTEM

Example: Show account infos
  $ fox ad -j NTDS.dit SYSTEM
`)

type Ad struct {
	// account flags
	Json  bool `short:"j" xor:"json,jsonl"`
	Jsonl bool `short:"J" xor:"json,jsonl"`

	// secrets flags
	History bool `short:"H"`
	Lm      bool
	Nt      bool

	// paths
	Paths []string `arg:"" optional:""`
}

func (cmd *Ad) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	if cmd.Lm || cmd.Nt {
		cmd.Json = false
		cmd.Jsonl = false
	}

	return nil
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

	pek, acc, err := extract.Extract(f1.Bytes(), key)

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

	for _, a := range acc {
		text.Match(cmd.format(record.New(a)), cli.Regexp)
	}

	if cli.Verbose > 1 {
		log.Printf("found %d account(s)\n", len(acc))
	}

	return nil
}

func (cmd *Ad) format(r *record.Record) string {
	switch {
	case cmd.Jsonl:
		return text.ColorizeAs(r.ToJSONL(), "json")
	case cmd.Json:
		return text.ColorizeAs(r.ToJSON(), "json")
	case cmd.Lm && cmd.Nt:
		return fmt.Sprintf("%s:%s", r.LMHash, r.NTHash)
	case cmd.Lm:
		return r.LMHash
	case cmd.Nt:
		return r.NTHash
	default:
		return r.ToNTLM(cmd.History)
	}
}
