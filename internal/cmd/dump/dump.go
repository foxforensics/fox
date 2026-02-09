package dump

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/cuhsat/go-secretsdump/pkg/ntds"
	//"github.com/mxk/go-vss"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
)

type Hash struct {
	Nt string `json:"nt,omitempty"`
	Lm string `json:"lm,omitempty"`
}

var Usage = strings.TrimSpace(`
Dumps sensitive data.

fox dump [FLAGS...] system ntds.dit

Flags:
  -V, --vss                dumps data using a Volume Shadow Copy (VSS)
  -j, --json               dumps data as JSON objects
  -J, --jsonl              dumps data as JSON lines
  -N, --nt                 shows only the NT hashes
  -L, --lm                 shows only the LM hashes

Examples:
  $ fox dump system ntds.dit
`)

const agree = "I understand"

type Dump struct {
	Vss   bool `short:"V"`
	Yes   bool `hidden:""`
	Json  bool `short:"j" xor:"json,jsonl"`
	Jsonl bool `short:"J" xor:"json,jsonl"`
	Nt    bool `long:"nt" xor:"nt,lm"` // hashcat -m 1000 / john --format=NT
	Lm    bool `long:"lm" xor:"nt,lm"` // hashcat -m 3000 / john --format=LM

	// paths
	Paths []string `arg:"" type:"path" optional:""`
}

func (cmd *Dump) Validate() error {
	if cmd.Vss && !cmd.Yes {
		log.Println(text.Warn("USING VSS WILL ALTER THE FILESYSTEM!!"))
		log.Println(text.Warn(fmt.Sprintf("PLEASE TYPE '%s' TO PROCEED", agree)))

		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		if input.Text() != agree {
			return errors.New("aborted")
		}
	}

	return nil
}

func (cmd *Dump) Run(cli *cli.Globals) error {
	if len(cmd.Paths) != 2 {
		fmt.Println(Usage)
		return nil
	}

	cli.NoConvert = true // forced

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	reg := <-ch
	dit := <-ch

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

	reg.Discard()
	dit.Discard()

	return nil
}

func (cmd *Dump) format(c *ntds.Credentials, re *regexp.Regexp) string {
	var line string
	var data any = c

	switch {
	case cmd.Nt:
		data = Hash{Nt: c.Nt}
	case cmd.Lm:
		data = Hash{Lm: c.Lm}
	}

	switch {
	case cmd.Jsonl:
		b, _ := json.MarshalIndent(data, "", "  ")
		line = text.ColorizeStringAs(string(b), "json")
	case cmd.Json:
		b, _ := json.Marshal(data)
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

/*
func (cmd *Dump) shadow(drv string) (dir string, err error) {
	dir = filepath.Join(os.TempDir(), "fox")

	if _, err = os.Stat(dir); !os.IsNotExist(err) {
		return "", errors.New("directory already exists")
	}

	err = vss.CreateLink(dir, drv)
}
*/
