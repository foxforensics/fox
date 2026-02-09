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
	registry "github.com/cuhsat/fox/v4/internal/pkg/data/extract/reg"

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

type Hash struct {
	Nt string `json:"nt,omitempty"`
	Lm string `json:"lm,omitempty"`
}

type Dump struct {
	Json  bool `short:"j" xor:"json,jsonl"`
	Jsonl bool `short:"J" xor:"json,jsonl"`

	// registry flags
	Bootkey bool `short:"K"`

	// active directory flags
	Nt bool `short:"N" long:"nt" xor:"nt,lm"` // hashcat -m 1000 / john --format=NT
	Lm bool `short:"L" long:"lm" xor:"nt,lm"` // hashcat -m 3000 / john --format=LM

	// paths
	Paths []string `arg:"" type:"path" optional:""`
}

func (cmd *Hash) String() string {
	switch {
	case len(cmd.Nt) > 0:
		return cmd.Nt
	case len(cmd.Lm) > 0:
		return cmd.Lm
	default:
		return "" // error
	}
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
		key, err := registry.BootKey(reg.Reader())

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

func (cmd *Dump) convert(c *ntds.Credentials) any {
	switch {
	case cmd.Nt:
		return Hash{Nt: c.Nt}
	case cmd.Lm:
		return Hash{Lm: c.Lm}
	default:
		return c
	}
}

func (cmd *Dump) format(v any, re *regexp.Regexp) string {
	var line string

	switch {
	case cmd.Jsonl:
		b, _ := json.MarshalIndent(v, "", "  ")
		line = text.ColorizeStringAs(string(b), "json")
	case cmd.Json:
		b, _ := json.Marshal(v)
		line = text.ColorizeStringAs(string(b), "json")
	case cmd.Nt:
		line = fmt.Sprint(v)
	case cmd.Lm:
		line = fmt.Sprint(v)
	default:
		line = fmt.Sprint(v)
	}

	if re != nil {
		line = text.MarkMatch(line, re)
	}

	return line
}
