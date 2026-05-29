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
fox hash [FLAGS...] <PATHS...>

Flags:
  -H, --hash=NAME,...      Show a specific hash (default: SHA256)
  -a, --all                Show all hashes and checksums
  -j, --json               Show results as JSON objects
  -J, --jsonl              Show results as JSON lines

Remarks:
  If more than one algorithm is specified, results will be grouped by path.

Example: Hash archive contents as MD5
  $ fox hash -Hmd5 files.7z

Example: Hash binaries for similarity
  $ fox hash -Himpfuzzy *.exe

Example: Hash binary inside an archive
  $ fox hash -Pinfected ioc.zip:ioc.exe

Cryptographic hashes (BLAKE family):
  BLAKE2S-256, BLAKE2B-256, BLAKE2B-384, BLAKE2B-512, BLAKE3-256, BLAKE3-512

Cryptographic hashes (GOST family):
  GOST2012-256, GOST2012-512

Cryptographic hashes (SHA family):
  SHA1, SHA224, SHA256, SHA512, SHA3, SHA3-224, SHA3-256, SHA3-384, SHA3-512

Cryptographic hashes (SKEIN family):
  SKEIN-224, SKEIN-256, SKEIN-384, SKEIN-512

Cryptographic hashes (MD family):
  MD2, MD4, MD5, MD6

Cryptographic hashes (other):
  HAS-160, LSH-256, LSH-512, RIPEMD-160, SHAKE128, SHAKE256, SM3, WHIRLPOOL

Performance hashes:
  DJB2, FNV-1, FNV-1A, MURMUR3, RAPIDHASH, SIPHASH, XXH32, XXH64, XXH3

Perceptual hashes:
  AVERAGE, DIFFERENCE, MEDIAN, PHASH, WHASH, MARRHILDRETH, BLOCKMEAN, PDQ, RASH

Similarity hashes:
  IMPFUZZY, IMPHASHO, IMPHASHS, SSDEEP, TLSH

Windows hashes:
  LM, NT, PE

Checksums:
  ADLER32, FLETCHER4, CRC16-CCITT, CRC32-C, CRC32-IEEE, CRC64-ECMA, CRC64-ISO
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
