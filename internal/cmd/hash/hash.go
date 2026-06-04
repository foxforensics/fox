package hash

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/alecthomas/kong"
	"go.foxforensics.dev/hasher/hash"
	"go.foxforensics.dev/rhash/database"

	cli "go.foxforensics.dev/fox/v4/internal/cmd"

	"go.foxforensics.dev/fox/v4/internal/pkg/tables"
	"go.foxforensics.dev/fox/v4/internal/pkg/text"
)

var Usage = strings.TrimSpace(`
Usage: fox hash [FLAGS...] <list|PATHS...>

Flags:
  -H, --hash=NAME,...      Use hash algorithm(s) (default: SHA256)
  -a, --all                Show all hashes and checksums
  -j, --json               Show results as JSON objects
  -J, --jsonl              Show results as JSON lines

Reverse flags:
  -l, --lookup             Lookup hash using wordlist
  -g, --guess              Guess the used algorithm(s)

Remarks:
  If 'list' is specified as path, only the built-in algorithms will be shown.
  If more than one algorithm is specified, results will be grouped by path.

Example: Hash archive contents as MD5
  $ fox hash -Hmd5 files.7z

Example: Hash binaries for similarity
  $ fox hash -Himpfuzzy *.exe

Example: Hash binary inside an archive
  $ fox hash -Pinfected ioc.zip:ioc.exe

Example: Lookup hash sum in wordlist
  $ fox hash -Hsha1 -l dump.sha1

Report bugs at: foxforensics.dev/issues
`)

type FileHash struct {
	File string            `json:"file,omitempty"`
	Hash map[string]string `json:"hash,omitempty"`
}

func (fh *FileHash) ToJSON() string {
	b, _ := json.MarshalIndent(fh, "", "  ")
	return string(b)
}

func (fh *FileHash) ToJSONL() string {
	b, _ := json.Marshal(fh)
	return string(b)
}

type Hash struct {
	Hash  []string `short:"H" sep:","`
	All   bool     `short:"a"`
	Json  bool     `short:"j" xor:"json,jsonl,lookup,guess"`
	Jsonl bool     `short:"J" xor:"json,jsonl,lookup,guess"`

	// reverse flags
	Lookup bool `short:"l" xor:"json,jsonl,lookup,guess"`
	Guess  bool `short:"g" xor:"json,jsonl,lookup,guess"`

	// paths
	Paths []string `arg:"" optional:""`
}

func (cmd *Hash) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	var err error

	if cmd.All || slices.Contains(cmd.Hash, "*") {
		cmd.Hash = list(!cmd.Lookup) // use all
	}

	// default algorithm
	if len(cmd.Hash) == 0 {
		cmd.Hash = []string{hash.SHA256}
	}

	// build tables
	if cmd.Lookup {
		_, err = tables.Build(nil, cmd.Hash...)
	}

	return err
}

func (cmd *Hash) Run(cli *cli.Globals) error {
	cmd.Paths = append(cmd.Paths, cli.Input...)

	if len(cmd.Paths) == 0 {
		return text.Usage(Usage)
	}

	if cmd.Paths[0] == "list" {
		for _, s := range list(true) {
			text.Write(s)
		}

		// exit early
		return nil
	}

	if !cli.NoPretty && len(cmd.Hash) == 1 {
		text.Title(cmd.Paths...)
	}

	plain, n := !cmd.Json && !cmd.Jsonl, 0

	for _, algo := range cmd.Hash {
		n = max(n, len(algo))
	}

	ch := cli.Load(cmd.Paths, true)
	defer cli.Discard()

	for h := range ch {
		fh := &FileHash{
			File: h.String(),
			Hash: make(map[string]string),
		}

		if !cli.NoPretty && len(cmd.Hash) > 1 {
			text.Title(h.String())
		}

		for _, k := range cmd.Hash {
			if cmd.Lookup {
				a, v := tables.Lookup(string(h.Bytes()))
				text.Match(fmt.Sprintf("%s  %s", v, strings.ToUpper(a)), cli.Regexp)
				break
			} else if cmd.Guess {
				for _, a := range collect(database.Lookup(string(h.Bytes()))) {
					text.Match(a, cli.Regexp)
				}
				continue
			}

			sum, err := hash.Sum(k, h.Bytes())

			if errors.Is(err, hash.NotSupported) {
				return fmt.Errorf("%s: %s", err, k)
			}

			if cli.Regexp != nil {
				if ok, _ := cli.Regexp.MatchString(sum); !ok {
					continue // do not include hash
				}
			}

			if err == nil {
				fh.Hash[k] = sum
			} else {
				fh.Hash[k] = err.Error()
			}

			if !plain {
				continue // will be formated
			}

			if err != nil {
				sum = text.AsGray(err.Error())
			}

			if len(cmd.Hash) > 1 {
				sum = fmt.Sprintf("%-*s  %s", n, strings.ToUpper(k), sum)
			} else {
				sum = fmt.Sprintf("%s  %s", sum, fh.File)
			}

			text.Match(sum, cli.Regexp)
		}

		if !plain {
			text.Match(cmd.format(fh), cli.Regexp)
		}

		h.Discard()
	}

	return nil
}

func (cmd *Hash) format(fh *FileHash) string {
	switch {
	case cmd.Jsonl:
		return text.ColorizeAs(fh.ToJSONL(), "json")
	case cmd.Json:
		return text.ColorizeAs(fh.ToJSON(), "json")
	default:
		return ""
	}
}

func collect(ch <-chan string) []string {
	v := make([]string, len(hash.Algorithms))

	for s := range ch {
		v = append(v, strings.Split(s, "\n")...)
	}

	return v
}

func list(all bool) []string {
	v := make([]string, len(hash.Algorithms))

	for _, a := range hash.Algorithms {
		if all || a.Type == hash.Cryptographic {
			v = append(v, a.Name)
		}
	}

	return v
}
