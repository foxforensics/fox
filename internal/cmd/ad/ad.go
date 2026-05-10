package ad

import (
	"log"
	"strings"

	"github.com/alecthomas/kong"
	"go.foxforensics.dev/bootkey/bootkey"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/rainbow"
	"go.foxforensics.dev/hashdump/extract"

	cli "go.foxforensics.dev/fox/v4/internal/cmd"

	"go.foxforensics.dev/fox/v4/internal/pkg/text"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/record"
)

var Usage = strings.TrimSpace(`
fox ad [FLAGS...] NTDS SYSTEM

Flags:
  -j, --json               Show accounts as JSON objects
  -J, --jsonl              Show accounts as JSON lines

Secrets flags:
  -L, --lookup             Lookup hashes in the rainbow table
  -H, --history            Extract also the users hash history
      --only-lm            Extract only the LM hashes (hashcat: 3000)
      --only-nt            Extract only the NT hashes (hashcat: 1000)

Example: Show NTLM secrets
  $ fox ad -LH NTDS.dit SYSTEM

Example: Show account infos
  $ fox ad -j NTDS.dit SYSTEM
`)

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
	if cmd.OnlyLm || cmd.OnlyNt {
		cmd.Json = false
		cmd.Jsonl = false
	}

	return nil
}

func (cmd *Ad) Run(cli *cli.Globals) error {
	if len(cmd.Paths) < 2 {
		return text.Usage(Usage)
	}

	if cmd.Lookup {
		if cli.Verbose > 1 {
			log.Println("building rainbow table")
		}

		err := rainbow.Build(cli.Parallel)

		if err != nil {
			return err
		}
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
	case cmd.OnlyLm:
		return r.OnlyLM()
	case cmd.OnlyNt:
		return r.OnlyNT()
	default:
		return r.ToNTLM(cmd.History)
	}
}
