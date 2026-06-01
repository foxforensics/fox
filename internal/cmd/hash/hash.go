package hash

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/alecthomas/kong"
	"go.foxforensics.dev/hasher/hash"

	cli "go.foxforensics.dev/fox/v4/internal/cmd"

	"go.foxforensics.dev/fox/v4/internal/pkg/text"
)

var Usage = strings.TrimSpace(`
Usage: fox hash [FLAGS...] <list|PATHS...>

Flags:
  -H, --hash=NAME,...      Use hash algorithm(s) (default: SHA256)
  -a, --all                Show all hashes and checksums
  -j, --json               Show results as JSON objects
  -J, --jsonl              Show results as JSON lines

Remarks:
  If 'list' is specified as path, only the built-in algorithms will be shown.
  If more than one algorithm is specified, results will be grouped by path.

Example: Hash archive contents as MD5
  $ fox hash -Hmd5 files.7z

Example: Hash binaries for similarity
  $ fox hash -Himpfuzzy *.exe

Example: Hash binary inside an archive
  $ fox hash -Pinfected ioc.zip:ioc.exe

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
	Json  bool     `short:"j" xor:"json,jsonl"`
	Jsonl bool     `short:"J" xor:"json,jsonl"`

	// paths
	Paths []string `arg:"" optional:""`
}

func (cmd *Hash) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	if cmd.All || slices.Contains(cmd.Hash, "*") {
		cmd.Hash = hash.Algorithms // use all
	}

	// default algorithm
	if len(cmd.Hash) == 0 {
		cmd.Hash = []string{hash.SHA256}
	}

	return nil
}

func (cmd *Hash) Run(cli *cli.Globals) error {
	cmd.Paths = append(cmd.Paths, cli.Input...)

	if len(cmd.Paths) == 0 {
		return text.Usage(Usage)
	}

	if cmd.Paths[0] == "list" {
		for _, s := range hash.Algorithms {
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
