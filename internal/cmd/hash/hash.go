package hash

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/alecthomas/kong"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/hash"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
)

var Usage = strings.TrimSpace(`
Prints file hashes and checksums.

fox hash [FLAGS ...] <PATHS ...>

Flags:
  -a, --all                use all available algorithms
  -T, --type=ALGO,...      use algorithms (default SHA256)
  -F, --find=HASH,...      show only files that match

Examples:
  $ fox hash -Tmd5,sha1 files.7z

Cryptographic hashes (BLAKE family):
  BLAKE2S-256, BLAKE2B-256, BLAKE2B-384, BLAKE2B-512, BLAKE3-256, BLAKE3-512

Cryptographic hashes (SHA family):
  SHA1, SHA256, SHA3, SHA3-224, SHA3-256, SHA3-384, SHA3-512

Cryptographic hashes (MD family):
  MD2, MD4, MD5, MD6

Performance hashes:
  XXH64, XXH3

Similarity hashes:
  SSDEEP, TLSH

Windows hashes:
  LM, NT, PE

Checksums:
  ADLER32, CRC32-C, CRC32-IEEE, CRC64-ECMA, CRC64-ISO
`)

type Hash struct {
	All   bool     `short:"a"`
	Type  []string `short:"T" sep:"," default:"SHA256"`
	Find  []string `short:"F" sep:","`
	Paths []string `arg:"" type:"path" optional:""`
}

func (cmd *Hash) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	if cmd.All {
		cmd.Type = hash.Algorithms
	}

	return nil
}

func (cmd *Hash) Run(cli *cli.Globals) error {
	if cli.Help || len(cmd.Paths) == 0 {
		fmt.Print(Usage)
		return nil
	}

	hs := cli.Load(cmd.Paths)
	defer cli.Discard()

	for _, typ := range cmd.Type {
		if !hash.Secure(typ) && !cli.NoWarning {
			log.Printf("used algorithm %s is not cryptically secure!\n", typ)
		}

		if len(cmd.Type) > 1 {
			_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide(text.Header(strings.ToUpper(typ))))
		}

		for _, h := range hs.Get() {
			sum, err := hash.Sum(typ, h.MMap())

			if err != nil {
				log.Println(err)
				continue
			}

			if len(cmd.Find) == 0 || slices.Contains(cmd.Find, sum) {
				_, _ = fmt.Fprintf(cli.Stdout, "%s  %s\n", sum, text.Hide(h))
			}
		}
	}

	return nil
}
