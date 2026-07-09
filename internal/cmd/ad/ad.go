package ad

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"go.foxforensics.eu/bootkey/bootkey"
	"go.foxforensics.eu/fox/v4/internal/cmd"
	"go.foxforensics.eu/fox/v4/internal/pkg/ad/record"
	"go.foxforensics.eu/fox/v4/internal/pkg/ad/tables"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/fox/v4/library/binaries/bin/ese"
	"go.foxforensics.eu/fox/v4/library/binaries/bin/reg"
	"go.foxforensics.eu/fox/v4/library/formats"
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
	cmd.Paths = append(cmd.Paths, fox.Paths...)

	if len(cmd.Paths) < 2 {
		sys.Usage(Usage)
		return nil
	}

	if len(cmd.Paths) > 2 {
		slog.Warn("additional paths will be ignored")
	}

	ch, err := fox.Init(cmd.Paths, true)

	if err != nil {
		return err
	}

	ntds := <-ch

	if ntds == nil {
		return errors.New("required file(s) missing")
	}

	defer ntds.Free()

	if !ese.Detect(ntds.Bytes()) {
		return errors.New("invalid file format")
	}

	hive := <-ch

	if hive == nil {
		return errors.New("required file(s) missing")
	}

	defer hive.Free()

	if !reg.Detect(hive.Bytes()) {
		return errors.New("invalid file format")
	}

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
	var a []fmt.Stringer

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
		fox.Writer.Match(cmd.format(v), fox.Regexp)
	}

	return len(a), nil
}

func (cmd *Ad) format(a fmt.Stringer) string {
	if v, ok := a.(*record.Secret); ok {
		switch {
		case cmd.LmOnly:
			return v.LmOnly()
		case cmd.NtOnly:
			return v.NtOnly()
		case !cmd.Json && !cmd.Jsonl:
			return v.ToNTLM(cmd.History)
		}
	}

	return formats.Auto(a, cmd.Json, cmd.Jsonl)
}
