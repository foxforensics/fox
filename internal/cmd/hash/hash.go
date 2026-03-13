package hash

import (
	"strings"

	"github.com/alecthomas/kong"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/hash"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

var Usage = strings.TrimSpace(`
Show file hashes and checksums.

fox hash [FLAGS...] <PATHS...>

Flags:
  -A, --algo=NAME,...      Use algorithm(s) (default: SHA256)
  -a, --all                Use all algorithms

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

Windows specific:
  LM, NT, PE Checksum

Checksums:
  ADLER32, FLETCHER4, CRC16-CCITT, CRC32-C, CRC32-IEEE, CRC64-ECMA, CRC64-ISO
`)

type Hash struct {
	Algo  []string `short:"A" sep:","`
	All   bool     `short:"a"`
	Paths []string `arg:"" optional:""`
}

func (cmd *Hash) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	if cmd.All {
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

	cli.NoConvert = true // forced

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	for h := range ch {
		if !cli.NoPretty && len(cmd.Algo) > 1 {
			text.Framed(h.String())
		}

		for _, algo := range cmd.Algo {
			sum, err := hash.Sum(algo, h.Bytes())

			if cli.Regexp != nil && !cli.Regexp.MatchString(sum) {
				continue // not matched afterward
			}

			res := text.MarkMatch(sum, cli.Regexp)

			if err != nil {
				res = text.AsGray(err.Error())
			}

			if !cli.NoPretty && len(cmd.Algo) > 1 {
				text.Pretty("%-21s  %s", text.AsBold(strings.ToUpper(algo)), res)
			} else if len(cmd.Algo) > 1 {
				text.Writeln("%-21s  %s", text.AsBold(strings.ToUpper(algo)), res)
			} else {
				text.Writeln("%s  %s", res, text.AsBold(h.Name))
			}
		}

		h.Discard()
	}

	return nil
}
