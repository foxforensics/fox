package dump

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/cuhsat/go-secretsdump/pkg/ntds"

	cli "github.com/cuhsat/fox/v4/internal/cmd"
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

	// print bootkey and exit early
	if cmd.Bootkey {
		key, err := sys.BootKey(reg.Reader())

		if err == nil {
			_, _ = fmt.Fprintln(cli.Stdout, fmt.Sprintf("%x", key))
		}

		return err
	}

	dit := <-ch
	defer dit.Discard()

	dump, err := ntds.New(
		bytes.NewReader(reg.Bytes()),
		bytes.NewReader(dit.Bytes()),
		len(dit.Bytes()),
	)

	if err != nil {
		log.Fatalln(err)
	}

	for c := range dump {
		line := cmd.format(&c, cli.Regexp)

		if cli.Regexp != nil && !cli.Regexp.MatchString(line) {
			continue // not matched afterward
		}

		_, _ = fmt.Fprintln(cli.Stdout, line)
	}

	return nil
}

func (cmd *Dump) format(c *ntds.Credentials, re *regexp.Regexp) string {
	var line string

	switch {
	case cmd.Jsonl:
		b, _ := json.MarshalIndent(c, "", "  ")
		line = text.ColorizeStringAs(string(b), "json")
	case cmd.Json:
		b, _ := json.Marshal(c)
		line = text.ColorizeStringAs(string(b), "json")
	case cmd.Nt:
		line = c.Nt
	case cmd.Lm:
		line = c.Lm
	default:
		line = c.String()
	}

	if re != nil {
		line = text.MarkMatch(line, re)
	}

	return line
}
