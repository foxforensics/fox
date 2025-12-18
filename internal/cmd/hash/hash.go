package hash

import (
	"fmt"
	"log"
	"slices"
	"strings"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/hash"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
)

var Usage = strings.TrimSpace(`
Prints file hashes and checksums.

fox hash [FLAGS ...] <PATHS ...>

Flags:
  -a, --algo=ALGO,...      use algorithms (default SHA256)
  -F, --find=HASH,...      show only files that match

Cryptographic hashes:
  MD2, MD4, MD5, SHA1, SHA256, SHA3, SHA3-224, SHA3-256, SHA3-384, SHA3-512

Performance hashes:
  XXH64, XXH3

Similarity hashes:
  SSDEEP, TLSH

Windows hashes:
  LM, NT, PE

Checksums:
  CRC32-C, CRC32-IEEE, CRC64-ECMA, CRC64-ISO
`)

type Hash struct {
	Algo  []string `short:"a" sep:"," default:"SHA256"`
	Find  []string `short:"F" sep:","`
	Paths []string `arg:"" type:"path" optional:""`
}

func (cmd *Hash) Run(cli *cli.Globals) error {
	if cli.Help || len(cmd.Paths) == 0 {
		fmt.Print(Usage)
		return nil
	}

	hs := cli.Bootstrap(cmd.Paths)
	defer cli.ThrowAway()

	for _, algo := range cmd.Algo {
		if !hash.Secure(algo) {
			log.Printf("used algorithm %s is not cryptically secure!\n", algo)
		}

		if len(cmd.Algo) > 1 {
			_, _ = fmt.Fprintf(cli.Stdout, "%s\n", text.Hide(text.Header(strings.ToUpper(algo))))
		}

		for _, h := range hs.Get() {
			sum, err := hash.Sum(algo, h.MMap())

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
