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
  -u, --use=ALGORITHM,...  use algorithms (default: SHA256)
  -a, --all                all algorithms

Example:
  $ fox hash -uTLSH files.7z

Remark:
  Results will be grouped by path, if more than one algorithm is specified.

Cryptographic hashes (BLAKE family):
  BLAKE2S-256, BLAKE2B-256, BLAKE2B-384, BLAKE2B-512, BLAKE3-256, BLAKE3-512

Cryptographic hashes (SHA family):
  SHA1, SHA224, SHA256, SHA512, SHA3, SHA3-224, SHA3-256, SHA3-384, SHA3-512

Cryptographic hashes (MD family):
  MD2, MD4, MD5, MD6

Cryptographic hashes (other):
  RIPEMD-160, SHAKE128, SHAKE256

Performance hashes:
  FNV-1, FNV-1A, MURMUR3, XXH32, XXH64, XXH3

Similarity hashes:
  SSDEEP, TLSH

Windows specific:
  LM, NT, PE Checksum

Image specific:
  AHASH, DHASH, PHASH

Checksums:
  ADLER32, CRC32-C, CRC32-IEEE, CRC64-ECMA, CRC64-ISO
`)

type Hash struct {
	Use   []string `short:"u" sep:"," default:"SHA256"`
	All   bool     `short:"a"`
	Paths []string `arg:"" type:"path" optional:""`
}

func (cmd *Hash) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	if cmd.All {
		cmd.Use = hash.Algorithms
	}

	return nil
}

func (cmd *Hash) Run(cli *cli.Globals) error {
	if cli.Help || len(cmd.Paths) == 0 {
		fmt.Print(Usage)
		return nil
	}

	cli.NoConvert = true // forced

	// compatibility mode
	if len(cmd.Use) == 1 {
		cli.NoFile = true
	}

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	for h := range ch {
		if !cli.NoFile {
			_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide(text.Header(h.String())))
		}

		for _, algo := range cmd.Use {
			if !hash.IsSecure(algo) && !cli.NoWarnings {
				log.Printf("warning: %s is not a cryptographic secure algorithm!\n", algo)
			}

			sum, err := hash.Sum(algo, h.MMap())

			if err != nil {
				log.Println(err)
				continue
			}

			if cli.Filter != nil && !cli.Filter.MatchString(sum) {
				continue // not matched afterward
			}

			sum = text.MarkMatch(sum, cli.Filter)

			if len(cmd.Use) > 1 {
				_, _ = fmt.Fprintf(cli.Stdout, "%s  %s\n", sum, text.Hide(strings.ToUpper(algo)))
			} else {
				_, _ = fmt.Fprintf(cli.Stdout, "%s  %s\n", sum, text.Hide(h.Name))
			}
		}

		h.Discard()
	}
	return nil
}
