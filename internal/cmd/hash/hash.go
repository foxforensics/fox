package hash

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/alecthomas/kong"
	"go.foxforensics.eu/fox/v4/internal/cmd"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/formats"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/tables"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/fox/v4/internal/sys/writer"
	"go.foxforensics.eu/hasher/hash"
	"go.foxforensics.eu/rhash/database"
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
  $ fox hash -Pinfected ioc.zip::ioc.exe

Example: Guess hash algorithm from sum
  $ fox hash -Hsha1 -g sum.txt

Report bugs at: foxforensics.eu/issues
`)

type FileHash struct {
	File string            `json:"file,omitempty"`
	Hash map[string]string `json:"hash,omitempty"`
}

func (fh *FileHash) String() string {
	return fh.File
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
	if cmd.All || slices.Contains(cmd.Hash, "*") {
		cmd.Hash = list(!cmd.Lookup) // use all
	}

	// default algorithm
	if len(cmd.Hash) == 0 {
		cmd.Hash = []string{hash.SHA256}
	}

	return nil
}

func (cmd *Hash) Run(fox *cmd.Globals) error {
	cmd.Paths = append(cmd.Paths, fox.Paths...)

	if len(cmd.Paths) == 0 {
		return sys.Usage(Usage)
	}

	if cmd.Paths[0] == "list" {
		for _, s := range list(true) {
			fox.Writer.Write(s)
		}

		// exit early
		return nil
	}

	if cmd.Lookup {
		_, err := tables.Build(nil, fox.Threads, cmd.Hash...)

		if err != nil {
			return err
		}
	}

	plain, n := !cmd.Json && !cmd.Jsonl, 0

	for _, algo := range cmd.Hash {
		n = max(n, len(algo))
	}

	ch, err := fox.Init(cmd.Paths, true)

	if err != nil {
		return err
	}

	for h := range ch {
		fh := &FileHash{
			File: h.String(),
			Hash: make(map[string]string),
		}

		if !fox.NoPretty && len(cmd.Hash) > 1 {
			fox.Writer.Header(h.String())
		}

		for _, k := range cmd.Hash {
			if cmd.Lookup {
				a, v := tables.Lookup(string(h.Bytes()))
				fox.Writer.Match(fmt.Sprintf("%s  %s", v, strings.ToUpper(a)), fox.Regexp)
				break
			} else if cmd.Guess {
				for _, a := range collect(database.Lookup(fox.Context, string(h.Bytes()))) {
					fox.Writer.Match(a, fox.Regexp)
				}
				continue
			}

			sum, err := hash.Sum(k, h.Bytes())

			if errors.Is(err, hash.NotSupported) {
				return fmt.Errorf("%s: %s", err, k)
			}

			if fox.Regexp != nil && !plain {
				if ok, _ := fox.Regexp.MatchString(sum); !ok {
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

			pad := strings.Repeat(" ", n-len(k))
			typ := writer.AsGray(strings.ToUpper(k))

			if err != nil {
				sum = writer.AsGray(err.Error())
			}

			if len(cmd.Hash) > 1 {
				sum = fmt.Sprintf("%s%s  %s", typ, pad, sum)
			} else {
				sum = fmt.Sprintf("%s  %s", sum, fh.File)
			}

			fox.Writer.Match(sum, fox.Regexp)
		}

		if !plain {
			fox.Writer.Match(formats.Auto(fh, cmd.Json, cmd.Jsonl), fox.Regexp)
		}

		h.Discard()
	}

	return nil
}

func collect(ch <-chan string) []string {
	v := make([]string, 0, len(hash.Algorithms))

	for s := range ch {
		v = append(v, strings.Split(s, "\n")...)
	}

	return v
}

func list(all bool) []string {
	v := make([]string, 0, len(hash.Algorithms))

	for _, a := range hash.Algorithms {
		if all || a.Type == hash.Cryptographic {
			v = append(v, a.Name)
		}
	}

	return v
}
