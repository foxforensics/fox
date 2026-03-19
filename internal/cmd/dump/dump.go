package dump

import (
	"log"
	"regexp"
	"strings"

	"github.com/alecthomas/kong"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/data/extract/dit"
	"github.com/cuhsat/fox/v4/internal/pkg/data/extract/reg"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
)

var Usage = strings.TrimSpace(`
Dump Active Directory secrets.

fox dump [FLAGS...] SYSTEM [NTDS]

Flags:
  -j, --json               Dump data as JSON objects
  -J, --jsonl              Dump data as JSON lines

Registry flags:
  -K, --bootkey            Dump the host bootkey

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
	if len(cmd.Paths) != 2 && (len(cmd.Paths) == 1 && !cmd.Bootkey) {
		return text.Usage(Usage)
	}

	ch := cli.LoadPlain(cmd.Paths)
	defer cli.Discard()

	f1 := <-ch
	defer f1.Discard()

	key, err := reg.BootKey(f1.Reader())

	if err != nil {
		return err
	}

	if cmd.Bootkey {
		text.Write("Bootkey %x", key)
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

	for _, rec := range res {
		line := cmd.format(&rec, cli.Regexp)

		if cli.Regexp != nil && !cli.Regexp.MatchString(line) {
			continue // not matched afterward
		}

		text.Write(line)
	}

	if cli.Verbose > 0 {
		log.Println("dump: finished")
	}

	if cli.Verbose > 1 {
		log.Printf("dump: found %d records(s)\n", len(res))
	}

	return nil
}

func (cmd *Dump) format(rec *dit.Record, re *regexp.Regexp) string {
	var line string

	switch {
	case cmd.Jsonl:
		line = text.ColorizeAs(rec.ToJSONL(), "json")
	case cmd.Json:
		line = text.ColorizeAs(rec.ToJSON(), "json")
	case cmd.OnlyNt:
		line = rec.NtHash
	case cmd.OnlyLm:
		line = rec.LmHash
	default:
		line = rec.String()
	}

	if re != nil {
		line = text.MarkMatch(line, re)
	}

	return line
}
