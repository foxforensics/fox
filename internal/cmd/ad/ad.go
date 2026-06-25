package ad

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"go.foxforensics.eu/bootkey/bootkey"
	"go.foxforensics.eu/fox/v4/internal/cmd"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/binary/bin/ese"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/binary/bin/reg"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/format"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/record"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/tables"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/fox/v4/internal/sys/terminal"
	"go.foxforensics.eu/hashdump/extract"
	"go.foxforensics.eu/hasher/hash"
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

func (cmd *Ad) Run(fox *cmd.Globals) error {
	cmd.Paths = append(cmd.Paths, fox.Input...)

	if len(cmd.Paths) < 2 {
		return sys.Usage(Usage)
	}

	// paths will be loaded in order
	ch, err := fox.Init(cmd.Paths, true)

	if err != nil {
		return err
	}

	defer fox.Discard()

	ntds := <-ch

	if ntds == nil {
		return errors.New("required file(s) missing")
	}

	if !ese.Detect(ntds.Bytes()) {
		return errors.New("invalid file format")
	}

	defer ntds.Discard()

	hive := <-ch

	if hive == nil {
		return errors.New("required file(s) missing")
	}

	if !reg.Detect(hive.Bytes()) {
		return errors.New("invalid file format")
	}

	defer hive.Discard()

	if cmd.Lookup {
		slog.Info("building tables")

		n, err := tables.Build(cmd.Wordlist, fox.Threads, hash.LM, hash.NT)

		if err != nil {
			return err
		}

		slog.Debug(fmt.Sprintf("using %d NTLM hashes", n))
	}

	key, err := bootkey.ExtractFromReader(hive.Reader())

	if err != nil {
		return err
	}

	slog.Debug(fmt.Sprintf("BootKey %x", key))

	pek, err := extract.Keys(fox.Context, ntds.Bytes(), key)

	if err != nil {
		return err
	}

	for i, k := range pek {
		slog.Debug(fmt.Sprintf("PEK #%d %x", i, k))
	}

	n, err := cmd.extract(fox, key, ntds.Bytes())

	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("found %d records(s)", n))

	return nil
}

func (cmd *Ad) extract(fox *cmd.Globals, k, b []byte) (int, error) {
	var a []any

	switch {
	case cmd.Users:
		if v, err := extract.Accounts(fox.Context, b, k); err != nil {
			return 0, err
		} else {
			for _, r := range v {
				a = append(a, &record.User{Account: r})
			}
		}

	case cmd.Groups:
		if v, err := extract.Groups(fox.Context, b); err != nil {
			return 0, err
		} else {
			for _, r := range v {
				a = append(a, &record.Group{Group: r})
			}
		}

	case cmd.Computers:
		if v, err := extract.Computers(fox.Context, b); err != nil {
			return 0, err
		} else {
			for _, r := range v {
				a = append(a, &record.Computer{Computer: r})
			}
		}

	default:
		if v, err := extract.Accounts(fox.Context, b, k); err != nil {
			return 0, err
		} else {
			for _, r := range v {
				a = append(a, &record.Secret{Account: r})
			}
		}
	}

	for _, v := range a {
		sys.Stdout.Match(cmd.format(v), fox.Regexp)
	}

	return len(a), nil
}

func (cmd *Ad) format(a any) string {
	switch v := a.(type) {
	case *record.Secret:
		switch {
		case cmd.LmOnly:
			return v.LmOnly()
		case cmd.NtOnly:
			return v.NtOnly()
		default:
			return v.ToNTLM(cmd.History)
		}

	case record.Record:
		switch {
		case cmd.Jsonl:
			return terminal.ColorizeAs(format.AsJSONL(v), "json")
		case cmd.Json:
			return terminal.ColorizeAs(format.AsJSON(v), "json")
		default:
			return v.String()
		}
	}

	return ""
}
