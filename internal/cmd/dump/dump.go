package dump

import (
	"log"
	"strings"

	"github.com/alecthomas/kong"
	"go.foxforensics.dev/bootkey/bootkey"
	"go.foxforensics.dev/hashdump/hashdump"

	cli "go.foxforensics.dev/fox/v4/internal/cmd"

	"go.foxforensics.dev/fox/v4/internal/pkg/text"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/record"
)

var Usage = strings.TrimSpace(`
fox dump [FLAGS...] SYSTEM [NTDS]

Flags:
  -j, --json               Dump data as JSON objects
  -J, --jsonl              Dump data as JSON lines

Registry flags:
  -K, --bootkey            Dump the host bootkey

Active Directory flags:
      --only-lm            Extract only the LM hashes (hashcat: 3000)
      --only-nt            Extract only the NT hashes (hashcat: 1000)

Example: Dump the BootKey from registry
  $ fox dump system -K

Example: Dump NTLM password hashes
  $ fox dump system ntds.dit
`)

type Dump struct {
	Json  bool `short:"j" xor:"json,jsonl"`
	Jsonl bool `short:"J" xor:"json,jsonl"`

	// registry flags
	Bootkey bool `short:"K"`

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
	if (cmd.Bootkey && len(cmd.Paths) < 1) || (!cmd.Bootkey && len(cmd.Paths) < 2) {
		return text.Usage(Usage)
	}

	ch := cli.Load(cmd.Paths, true)
	defer cli.Discard()

	f1 := <-ch
	defer f1.Discard()

	key, err := bootkey.ReadData(f1.Reader())

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

	if cli.Verbose > 1 {
		log.Printf("dump: BootKey %x\n", key)
	}

	rec, pek, err := hashdump.Dump(f2.Bytes(), key)

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

	for _, r := range rec {
		text.Match(cmd.format(record.New(r)), cli.Regexp)
	}

	if cli.Verbose > 0 {
		log.Println("dump: finished")
	}

	if cli.Verbose > 1 {
		log.Printf("dump: found %d records(s)\n", len(rec))
	}

	return nil
}

func (cmd *Dump) format(r *record.Record) string {
	switch {
	case cmd.Jsonl:
		return text.ColorizeAs(r.ToJSONL(), "json")
	case cmd.Json:
		return text.ColorizeAs(r.ToJSON(), "json")
	case cmd.OnlyNt:
		return r.NT
	case cmd.OnlyLm:
		return r.LM
	default:
		return r.String()
	}
}
