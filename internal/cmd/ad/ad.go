package ad

import (
	"errors"
	"log"
	"strings"

	"go.foxforensics.dev/bootkey/bootkey"
	"go.foxforensics.dev/hashdump/extract"

	cli "go.foxforensics.dev/fox/v4/internal/cmd"

	"go.foxforensics.dev/fox/v4/internal/pkg/tables"
	"go.foxforensics.dev/fox/v4/internal/pkg/text"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/record"
)

var Usage = strings.TrimSpace(`
fox ad [FLAGS...] NTDS SYSTEM

Flags:
  -j, --json               Show AD records as JSON objects
  -J, --jsonl              Show AD records as JSON lines

Record flags:
  -U, --users              Extract all user records
  -G, --groups             Extract all group records
  -C, --computers          Extract all computer records

Secret flags:
  -L, --lookup             Lookup hashes with rainbow tables
  -H, --history            Extract also the users hash history
      --only-lm            Extract only the LM hashes (Hashcat mode 3000)
      --only-nt            Extract only the NT hashes (Hashcat mode 1000)

Remarks:
  If no records are specified, hashes will be shown in secretsdump manner.

Example: Show user records
  $ fox ad -jU NTDS.dit SYSTEM

Example: Show NTLM hashes
  $ fox ad -LH NTDS.dit SYSTEM
`)

type Ad struct {
	Json  bool `short:"j" xor:"json,jsonl"`
	Jsonl bool `short:"J" xor:"json,jsonl"`

	// record flags
	Users     bool `short:"U" xor:"users,groups,computers"`
	Groups    bool `short:"G" xor:"users,groups,computers"`
	Computers bool `short:"C" xor:"users,groups,computers"`

	// secret flags
	Lookup  bool `short:"L"`
	History bool `short:"H"`
	OnlyLm  bool `long:"only-lm" xor:"only-lm,only-nt"`
	OnlyNt  bool `long:"only-nt" xor:"only-lm,only-nt"`

	// hidden
	Wordlist []byte `hidden:"" type:"filecontent"`

	// paths
	Paths []string `arg:"" optional:""`
}

func (cmd *Ad) Run(cli *cli.Globals) error {
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

		n, err := tables.Build(cmd.Wordlist)

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

	if !cli.NoPretty {
		text.Title(ntds.String())
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
		text.Match(cmd.format(v), cli.Regexp)
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
		case cmd.OnlyLm:
			return v.OnlyLM()
		case cmd.OnlyNt:
			return v.OnlyNT()
		default:
			return v.ToNTLM(cmd.History)
		}
	}

	return ""
}
