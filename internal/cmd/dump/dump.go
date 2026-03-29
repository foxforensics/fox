package dump

import (
	"log"
	"strings"

	"github.com/alecthomas/kong"

	cli "go.foxforensics.dev/fox/v4/internal/cmd"

	"go.foxforensics.dev/fox/v4/internal/pkg/file/extract/dit"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/extract/reg"
	"go.foxforensics.dev/fox/v4/internal/pkg/text"
)

var Usage = strings.TrimSpace(`
fox dump [FLAGS...] SYSTEM [NTDS]

Flags:
  -j, --json               Dump data as JSON objects
  -J, --jsonl              Dump data as JSON lines

Registry flags:
  -B, --bootkey            Dump the host bootkey

Active Directory flags:
      --only-lm            Extract only the LM hashes (hashcat: 3000)
      --only-nt            Extract only the NT hashes (hashcat: 1000)

Examples:
  $ fox dump system ntds.dit
`)

type Dump struct {
	Json  bool `short:"j" xor:"json,jsonl"`
	Jsonl bool `short:"J" xor:"json,jsonl"`

	// registry flags
	Bootkey bool `short:"B"`

	// active directory flags
	OnlyLm bool `long:"only-lm" xor:"only-lm,only-nt"`
	OnlyNt bool `long:"only-nt" xor:"only-lm,only-nt"`

	// paths
	Paths []string `arg:"" optional:""`
}

func (cmd *Dump) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	if cmd.OnlyLm || cmd.OnlyNt {
		cmd.Json = false
		cmd.Jsonl = false
	}

	return nil
}

func (cmd *Dump) Run(cli *cli.Globals) error {
	if len(cmd.Paths) != 2 && (len(cmd.Paths) == 1 && !cmd.Bootkey) {
		return text.Usage(Usage)
	}

	ch := cli.Load(cmd.Paths, true)
	defer cli.Discard()

	f1 := <-ch
	defer f1.Discard()

	key, err := reg.BootKey(f1.Reader())

	if err != nil {
		return err
	}

	if cmd.Bootkey {
		if !cli.NoPretty {
			text.Title(f1.String())
		}

		text.Write("%x", key)
		return nil
	}

	f2 := <-ch
	defer f2.Discard()

	if cli.Verbose > 0 {
		log.Println("dump: started")
	}

	res, pek, err := dit.Extract(f2.Bytes(), key)

	if err != nil {
		return err
	}

	if cli.Verbose > 1 {
		for i, k := range pek {
			log.Printf("dump: PEK #%d %x\n", i, k)
		}
	}

	if !cli.NoPretty {
		text.Title(f2.String())
	}

	for _, r := range res {
		text.Match(cmd.format(&r), cli.Regexp)
	}

	if cli.Verbose > 0 {
		log.Println("dump: finished")
	}

	if cli.Verbose > 1 {
		log.Printf("dump: found %d records(s)\n", len(res))
	}

	return nil
}

func (cmd *Dump) format(r *dit.Record) string {
	switch {
	case cmd.Jsonl:
		return text.ColorizeAs(r.ToJSONL(), "json")
	case cmd.Json:
		return text.ColorizeAs(r.ToJSON(), "json")
	case cmd.OnlyNt:
		return r.NtHash
	case cmd.OnlyLm:
		return r.LmHash
	default:
		return r.String()
	}
}
