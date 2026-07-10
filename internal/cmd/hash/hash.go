package hash

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/alecthomas/kong"
	"go.foxforensics.eu/fox/v5/internal/cmd"
	"go.foxforensics.eu/fox/v5/internal/pkg"
	"go.foxforensics.eu/fox/v5/internal/sys"
	"go.foxforensics.eu/fox/v5/internal/sys/writer"
	"go.foxforensics.eu/fox/v5/library/formats"
	"go.foxforensics.eu/hasher/hash"
)

var Usage = strings.TrimSpace(`
Usage: fox hash [FLAGS...] <list|PATHS...>

Flags:
  -H, --hash=NAME,...      Use hash algorithm(s) (default: SHA256)
  -a, --all                Show all hashes and checksums
  -j, --json               Show results as JSON objects
  -J, --jsonl              Show results as JSON lines

Filter flags
  -B, --include=FILE       Include only known bad hashes
  -G, --exclude=FILE       Exclude all known good hashes

Remarks:
  If 'list' is specified as path, only the built-in algorithms will be shown.
  If more than one algorithm is specified, results will be grouped by path.

Example: Hash archive contents as MD5
  $ fox hash -Hmd5 files.7z

Example: Hash binaries for similarity
  $ fox hash -Himpfuzzy *.exe

Example: Hash binary inside an archive
  $ fox hash -Pinfected ioc.zip::ioc.exe

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
	Hash  []string `short:"H" xor:"hash,all" sep:","`
	All   bool     `short:"a" xor:"hash,all"`
	Json  bool     `short:"j" xor:"json,jsonl"`
	Jsonl bool     `short:"J" xor:"json,jsonl"`

	// filter flags
	Include []byte `short:"B" xor:"include,exclude" type:"filecontent"`
	Exclude []byte `short:"G" xor:"include,exclude" type:"filecontent"`

	// paths
	Paths []string `arg:"" optional:""`

	// internal
	include map[string]pkg.Nil `kong:"-"`
	exclude map[string]pkg.Nil `kong:"-"`
}

func (cmd *Hash) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	if cmd.All {
		for _, a := range hash.Algorithms {
			cmd.Hash = append(cmd.Hash, a.Name)
		}
	}

	// default algorithm
	if len(cmd.Hash) == 0 {
		cmd.Hash = []string{hash.SHA256}
	}

	if len(cmd.Hash) > 1 && (len(cmd.Include)+len(cmd.Exclude) > 0) {
		return errors.New("filters can not be used with multiple algorithms")
	}

	if len(cmd.Include) > 0 {
		cmd.include = sys.ParseMap(cmd.Include)
	}

	if len(cmd.Exclude) > 0 {
		cmd.exclude = sys.ParseMap(cmd.Exclude)
	}

	return nil
}

func (cmd *Hash) Run(fox *cmd.Globals) error {
	cmd.Paths = append(cmd.Paths, fox.Paths...)

	if len(cmd.Paths) == 0 {
		sys.Usage(Usage)
		return nil
	}

	if cmd.Paths[0] == "list" {
		for _, a := range hash.Algorithms {
			fmt.Println(a.Name)
		}

		// exit early
		return nil
	}

	plain, n := !cmd.Json && !cmd.Jsonl, 0

	for _, algo := range cmd.Hash {
		n = max(n, len(algo))
	}

	heaps, err := fox.Init(cmd.Paths, true)

	if err != nil {
		return err
	}

	for h := range heaps {
		fh := &FileHash{
			File: h.String(),
			Hash: make(map[string]string),
		}

		if !fox.NoPretty && len(cmd.Hash) > 1 {
			fox.Writer.FileHeader(h.String())
		}

		for _, k := range cmd.Hash {
			select {
			case <-fox.Context.Done():
				h.Free()
				return nil // canceled

			default:
			}

			sum, err := hash.Sum(k, h.Bytes())

			if errors.Is(err, hash.ErrNotSupported) {
				slog.Error(fmt.Sprintf("%s: %s", err, k))
				continue // invalid hash
			}

			if len(sum) == 0 {
				if len(cmd.Exclude)+len(cmd.Include) == 0 {
					slog.Debug("hash was empty")
				} else {
					continue // empty sum
				}
			}

			if len(cmd.Exclude) > 0 {
				if _, ok := cmd.exclude[sum]; ok {
					slog.Debug(fmt.Sprintf("hash was excluded %s", sum))
					continue // was excluded
				}
			}

			if len(cmd.Include) > 0 {
				if _, ok := cmd.include[sum]; ok {
					slog.Debug(fmt.Sprintf("hash was included %s", sum))
				} else {
					continue // not included
				}
			}

			if fox.Regexp != nil && !plain {
				if ok, err := fox.Regexp.MatchString(sum); !ok {
					if err != nil {
						slog.Error(err.Error())
					}
					continue // hash does not match
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

		// show only valid entries
		if !plain && len(fh.Hash) > 0 {
			fox.Writer.Match(formats.Auto(fh, cmd.Json, cmd.Jsonl), fox.Regexp)
		}

		h.Free()
	}

	return nil
}
