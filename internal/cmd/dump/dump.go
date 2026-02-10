package dump

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/data/extract/dit"
	"github.com/cuhsat/fox/v4/internal/pkg/data/extract/reg"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
)

var Usage = strings.TrimSpace(`
Dumps sensitive data.

fox dump [FLAGS...] system [ntds.dit]

Flags:
  -j, --json               dumps data as JSON objects
  -J, --jsonl              dumps data as JSON lines

Registry flags:
  -K, --bootkey            extracts only the host bootkey

Active Directory flags:
  -L, --lm                 extracts only the LM hashes (hashcat: 3000)
  -N, --nt                 extracts only the NT hashes (hashcat: 1000)

Examples:
  $ fox dump system ntds.dit
`)

type Dump struct {
	Json  bool `short:"j" xor:"json,jsonl,lm,nt"`
	Jsonl bool `short:"J" xor:"json,jsonl,lm,nt"`

	// registry flags
	Bootkey bool `short:"K"`

	// active directory flags
	Lm bool `short:"L" long:"lm" xor:"json,jsonl,lm,nt"`
	Nt bool `short:"N" long:"nt" xor:"json,jsonl,lm,nt"`

	// paths
	Paths []string `arg:"" type:"path" optional:""`
}

func (cmd *Dump) Run(cli *cli.Globals) error {
	if len(cmd.Paths) < 2 || (len(cmd.Paths) < 1 && cmd.Bootkey) {
		fmt.Println(Usage)
		return nil
	}

	cli.NoConvert = true // forced

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	f1 := <-ch
	defer f1.Discard()

	key, err := reg.BootKey(f1.Reader())

	if err != nil {
		return err
	}

	if cmd.Bootkey {
		_, _ = fmt.Fprintln(cli.Stdout, fmt.Sprintf("%x", key))
		return nil
	}

	f2 := <-ch
	defer f2.Discard()

	if cli.Verbose > 0 {
		log.Println("dump: started")
	}

	res, err := dit.Extract(f2.Bytes(), key)

	if err != nil {
		return err
	}

	for _, rec := range res {
		line := cmd.format(&rec, cli.Regexp)

		if cli.Regexp != nil && !cli.Regexp.MatchString(line) {
			continue // not matched afterward
		}

		_, _ = fmt.Fprintln(cli.Stdout, line)
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
		line = text.ColorizeStringAs(rec.ToJSONL(), "json")
	case cmd.Json:
		line = text.ColorizeStringAs(rec.ToJSON(), "json")
	case cmd.Nt:
		line = rec.Nt
	case cmd.Lm:
		line = rec.Lm
	default:
		line = rec.String()
	}

	if re != nil {
		line = text.MarkMatch(line, re)
	}

	return line
}
