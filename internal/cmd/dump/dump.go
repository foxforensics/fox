package dump

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	cli "github.com/cuhsat/fox/v4/internal/cmd"
	ad "github.com/cuhsat/fox/v4/internal/pkg/data/extract/dit"
	sys "github.com/cuhsat/fox/v4/internal/pkg/data/extract/reg"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
)

var Usage = strings.TrimSpace(`
Dumps sensitive data.

fox dump [FLAGS...] system [ntds.dit]

Flags:
  -j, --json               dumps data as JSON objects
  -J, --jsonl              dumps data as JSON lines

Registry flags:
  -K, --bootkey            extracts only the bootkey

Active Directory flags:
  -N, --nt                 extracts only the NT hashes
  -L, --lm                 extracts only the LM hashes

Examples:
  $ fox dump system ntds.dit
`)

type Dump struct {
	Json  bool `short:"j" xor:"json,jsonl,nt,lm"`
	Jsonl bool `short:"J" xor:"json,jsonl,nt,lm"`

	// registry flags
	Bootkey bool `short:"K"`

	// active directory flags
	Nt bool `short:"N" long:"nt" xor:"json,jsonl,nt,lm"` // hashcat -m 1000 / john --format=NT
	Lm bool `short:"L" long:"lm" xor:"json,jsonl,nt,lm"` // hashcat -m 3000 / john --format=LM

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

	reg := <-ch
	defer reg.Discard()

	key, err := sys.BootKey(reg.Reader())

	if err != nil {
		return err
	}

	// print bootkey and exit early
	if cmd.Bootkey {
		_, _ = fmt.Fprintln(cli.Stdout, fmt.Sprintf("%x", key))
		return nil
	}

	dit := <-ch
	defer dit.Discard()

	rec, err := ad.Extract(dit.Bytes(), key)

	if err != nil {
		return err
	}

	for _, r := range rec {
		line := cmd.format(&r, cli.Regexp)

		if cli.Regexp != nil && !cli.Regexp.MatchString(line) {
			continue // not matched afterward
		}

		_, _ = fmt.Fprintln(cli.Stdout, line)
	}

	return nil
}

func (cmd *Dump) format(r *ad.Record, re *regexp.Regexp) string {
	var line string

	switch {
	case cmd.Jsonl:
		b, _ := json.MarshalIndent(r, "", "  ")
		line = text.ColorizeStringAs(string(b), "json")
	case cmd.Json:
		b, _ := json.Marshal(r)
		line = text.ColorizeStringAs(string(b), "json")
	case cmd.Nt:
		line = r.Nt
	case cmd.Lm:
		line = r.Lm
	default:
		line = r.String()
	}

	if re != nil {
		line = text.MarkMatch(line, re)
	}

	return line
}
