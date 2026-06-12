package ad

import (
	"errors"
	"log"
	"strings"

	"go.foxforensics.eu/bootkey/bootkey"
	"go.foxforensics.eu/hashdump/extract"
	"go.foxforensics.eu/hasher/hash"

	cli "go.foxforensics.eu/fox/v4/internal/cmd"

	"go.foxforensics.eu/fox/v4/internal/pkg/tables"
	"go.foxforensics.eu/fox/v4/internal/pkg/text"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/record"
)

var Usage = strings.TrimSpace(`
Usage: fox ad [FLAGS...] NTDS SYSTEM

Record flags:
  -u, --users              Show all user records
  -g, --groups             Show all group records
  -c, --computers          Show all computer records
  -j, --json               Show records as JSON objects
  -J, --jsonl              Show records as JSON lines

Secret flags:
  -l, --lookup             Lookup hashes using wordlist
  -h, --history            Extract also the users hash history
      --lm-only            Extract only the LM hashes (Hashcat mode 3000)
      --nt-only            Extract only the NT hashes (Hashcat mode 1000)

Remarks:
  If no records are specified, hashes will be shown in secretsdump manner.

Example: Show NTLM hashes
  $ fox ad -hl NTDS.dit SYSTEM

Example: Show user records
  $ fox ad -uj NTDS.dit SYSTEM

Report bugs at: foxforensics.eu/issues
`)

type Ad struct {
	// record flags
	Users     bool `short:"u" xor:"users,groups,computers"`
	Groups    bool `short:"g" xor:"users,groups,computers"`
	Computers bool `short:"c" xor:"users,groups,computers"`
	Json      bool `short:"j" xor:"json,jsonl"`
	Jsonl     bool `short:"J" xor:"json,jsonl"`

	// secret flags
	Lookup  bool `short:"l"`
	History bool `short:"h"`
	LmOnly  bool `long:"lm-only" xor:"lm-only,nt-only"`
	NtOnly  bool `long:"nt-only" xor:"lm-only,nt-only"`

	// hidden
	Wordlist []byte `hidden:"" type:"filecontent"`

	// paths
	Paths []string `arg:"" optional:""`
}

func (cmd *Ad) Run(cli *cli.Globals) error {
	cmd.Paths = append(cmd.Paths, cli.Input...)

	if len(cmd.Paths) < 2 {
		return text.Usage(Usage)
	}

	ch := cli.Load(cmd.Paths, true)
	defer cli.Discard()

	ntds := <-ch

	if ntds == nil {
		return errors.New("required file(s) missing")
	}

	defer ntds.Discard()

	hive := <-ch

	if hive == nil {
		return errors.New("required file(s) missing")
	}

	defer hive.Discard()

	if cmd.Lookup {
		if cli.Verbose > 0 {
			log.Println("building tables")
		}

		n, err := tables.Build(cmd.Wordlist, hash.LM, hash.NT)

		if err != nil {
			return err
		}

		if cli.Verbose > 0 {
			log.Printf("using %d NTLM hashes\n", n)
		}
	}

	key, err := bootkey.ReadData(hive.Reader())

	if err != nil {
		return err
	}

	if cli.Verbose > 1 {
		log.Printf("BootKey %x\n", key)

		pek, err := extract.Keys(ntds.Bytes(), key)

		if err != nil {
			return err
		}

		for i, k := range pek {
			log.Printf("PEK #%d %x\n", i, k)
		}
	}

	n, err := cmd.extract(cli, key, ntds.Bytes())

	if err != nil {
		return err
	}

	if cli.Verbose > 1 {
		log.Printf("found %d records(s)\n", n)
	}

	return nil
}

func (cmd *Ad) extract(cli *cli.Globals, k, b []byte) (int, error) {
	var a []any

	switch {
	case cmd.Users:
		if v, err := extract.Accounts(b, k); err != nil {
			return 0, err
		} else {
			for _, r := range v {
				a = append(a, &record.User{Account: r})
			}
		}

	case cmd.Groups:
		if v, err := extract.Groups(b); err != nil {
			return 0, err
		} else {
			for _, r := range v {
				a = append(a, &record.Group{Group: r})
			}
		}

	case cmd.Computers:
		if v, err := extract.Computers(b); err != nil {
			return 0, err
		} else {
			for _, r := range v {
				a = append(a, &record.Computer{Computer: r})
			}
		}

	default:
		if v, err := extract.Accounts(b, k); err != nil {
			return 0, err
		} else {
			for _, r := range v {
				a = append(a, &record.Secret{Account: r})
			}
		}
	}

	for _, v := range a {
		text.Stdout.Match(cmd.format(v), cli.Regexp)
	}

	return len(a), nil
}

func (cmd *Ad) format(a any) string {
	switch v := a.(type) {
	case record.Record:
		switch {
		case cmd.Jsonl:
			return text.ColorizeAs(v.ToJSONL(), "json")
		case cmd.Json:
			return text.ColorizeAs(v.ToJSON(), "json")
		default:
			return v.String()
		}

	case *record.Secret:
		switch {
		case cmd.LmOnly:
			return v.LmOnly()
		case cmd.NtOnly:
			return v.NtOnly()
		default:
			return v.ToNTLM(cmd.History)
		}
	}

	return ""
}
