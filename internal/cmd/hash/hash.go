package hash

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/alecthomas/kong"

	cli "go.foxforensics.dev/fox/v4/internal/cmd"

	"go.foxforensics.dev/fox/v4/internal/pkg/hash"
	"go.foxforensics.dev/fox/v4/internal/pkg/text"
	"go.foxforensics.dev/fox/v4/internal/pkg/types"
)

type FileHash struct {
	File string         `json:"file,omitempty"`
	Hash map[string]any `json:"hash,omitempty"`
}

func (fh *FileHash) String() string {
	var sb strings.Builder

	n := 0

	for k := range maps.Keys(fh.Hash) {
		n = max(n, len(k))
	}

	for _, k := range slices.Sorted(maps.Keys(fh.Hash)) {
		v := fh.Hash[k]

		// render errors
		if err, ok := v.(error); ok {
			v = text.AsGray(err.Error())
		}

		if len(fh.Hash) > 1 {
			sb.WriteString(fmt.Sprintf("%-*s  %s\n", n, strings.ToUpper(k), v))
		} else {
			sb.WriteString(fmt.Sprintf("%s  %s", v, fh.File))
		}
	}

	return strings.TrimSpace(sb.String())
}

func (fh *FileHash) ToJSON() string {
	b, _ := json.MarshalIndent(fh, "", "  ")
	return string(b)
}

func (fh *FileHash) ToJSONL() string {
	b, _ := json.Marshal(fh)
	return string(b)
}

var Usage = strings.TrimSpace(`
Show file hashes and checksums.

fox hash [FLAGS...] <PATHS...>

Flags:
  -A, --algo=NAME,...      Show a specific hash (default: SHA256)
  -a, --all                Show all hashes and checksums
  -j, --json               Show results as JSON objects
  -J, --jsonl              Show results as JSON lines

Examples:
  $ fox hash -Amd5 files.7z

Remarks:
  Results will be grouped by path, if more than one algorithm is specified.

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
  FNV-1, FNV-1A, MURMUR3, RAPIDHASH, SIPHASH, XXH32, XXH64, XXH3

Perceptual hashes:
  AVERAGE, DIFFERENCE, MEDIAN, PHASH, WHASH, MARRHILDRETH, BLOCKMEAN, PDQ, RASH

Similarity hashes:
  IMPFUZZY, IMPHASH, IMPHASH0, SSDEEP, TLSH

Windows algorithms:
  LM, NT, PE

Checksums:
  ADLER32, FLETCHER4, CRC16-CCITT, CRC32-C, CRC32-IEEE, CRC64-ECMA, CRC64-ISO
`)

type Hash struct {
	Algo  []string `short:"A" sep:","`
	All   bool     `short:"a"`
	Json  bool     `short:"j" xor:"json,jsonl"`
	Jsonl bool     `short:"J" xor:"json,jsonl"`

	// paths
	Paths []string `arg:"" optional:""`
}

func (cmd *Hash) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	if cmd.All || slices.Contains(cmd.Algo, "*") {
		cmd.Algo = hash.Algorithms
	}

	// default algorithm
	if len(cmd.Algo) == 0 {
		cmd.Algo = []string{types.SHA256}
	}

	return nil
}

func (cmd *Hash) Run(cli *cli.Globals) error {
	if len(cmd.Paths)+len(cli.Paths) == 0 {
		return text.Usage(Usage)
	}

	if !cli.NoPretty && len(cmd.Algo) == 1 {
		text.Title(cmd.Paths...)
	}

	ch := cli.Load(cmd.Paths, true)
	defer cli.Discard()

	for h := range ch {
		fh := &FileHash{
			File: h.String(),
			Hash: make(map[string]any),
		}

		if !cli.NoPretty && len(cmd.Algo) > 1 {
			text.Title(h.String())
		}

		for _, algo := range cmd.Algo {
			sum, err := hash.Sum(algo, h.Bytes())

			if errors.Is(err, hash.ErrNotSupported) {
				return err // not supported
			}

			if cli.Regexp != nil && !cli.Regexp.MatchString(sum) {
				continue // do not include hash
			}

			if err == nil {
				fh.Hash[algo] = sum
			} else {
				fh.Hash[algo] = err
			}
		}

		text.Match(cmd.format(fh), cli.Regexp)

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
		return fh.String()
	}
}
