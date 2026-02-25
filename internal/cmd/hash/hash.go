package hash

import (
	"fmt"
	"log"
	"strings"

	"github.com/alecthomas/kong"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/hash"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
)

var Usage = strings.TrimSpace(`
Prints file hashes and checksums.

fox hash [FLAGS...] <PATHS...>

Flags:
  -A, --algo=NAME,...      uses algorithm(s) (default: SHA256)
  -a, --all                uses all algorithms

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

Similarity hashes:
  IMPFUZZY, IMPHASH, IMPHASH0, SSDEEP, TLSH

Windows specific:
  LM, NT, PE Checksum

Image specific:
  AHASH, DHASH, PHASH

Checksums:
  ADLER32, FLETCHER4, CRC16-CCITT, CRC32-C, CRC32-IEEE, CRC64-ECMA, CRC64-ISO
`)

type Hash struct {
	Algo  []string `short:"A" sep:"," default:"SHA256"`
	All   bool     `short:"a"`
	Paths []string `arg:"" optional:""`
}

func (cmd *Hash) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	if cmd.All {
		cmd.Algo = hash.Algorithms
	}

	return nil
}

func (cmd *Hash) Run(cli *cli.Globals) error {
	if len(cmd.Paths)+len(cli.Paths) == 0 {
		fmt.Println(Usage)
		return nil
	}

	cli.NoConvert = true // forced

	// compatibility mode
	if len(cmd.Algo) == 1 {
		cli.NoFile = true
	}

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	for h := range ch {
		if !cli.NoFile {
			_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide(text.Title(h.String())))
		}

		for _, algo := range cmd.Algo {
			if !hash.IsSecure(algo) && !cli.NoWarnings {
				log.Printf("warning: %s is not a cryptographically secure algorithm!\n", algo)
			}

			sum, err := hash.Sum(algo, h.Bytes())

			if err != nil {
				log.Println(err)
				continue
			}

			if cli.Regexp != nil && !cli.Regexp.MatchString(sum) {
				continue // not matched afterward
			}

			sum = text.MarkMatch(sum, cli.Regexp)

			if len(cmd.Algo) > 1 {
				_, _ = fmt.Fprintf(cli.Stdout, "%s  %s\n", sum, text.Hide(strings.ToUpper(algo)))
			} else {
				_, _ = fmt.Fprintf(cli.Stdout, "%s  %s\n", sum, text.Hide(h.Name))
			}
		}

		h.Discard()
	}
	return nil
}
