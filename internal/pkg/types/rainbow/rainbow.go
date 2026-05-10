// Package rainbow source: https://github.com/danielmiessler/SecLists/blob/master/Passwords/Common-Credentials/probable-v2_top-12000.txt
package rainbow

import (
	"bufio"
	"bytes"

	_ "embed"

	"github.com/klauspost/compress/zstd"
	"go.foxforensics.dev/hasher/hash"
)

//go:embed rainbow.zst
var table []byte

var tableLm = make(map[string]string, 12645)
var tableNt = make(map[string]string, 12645)

func Build() error {
	r, err := zstd.NewReader(bytes.NewReader(table))

	if err != nil {
		return err
	}

	defer r.Close()

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		b := scanner.Bytes()

		tableLm[hash.MustSum(hash.LM, b)] = string(b)
		tableNt[hash.MustSum(hash.NT, b)] = string(b)
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}

func Lookup(sum string) string {
	if v, ok := tableNt[sum]; ok {
		return v
	}

	if v, ok := tableLm[sum]; ok {
		return v
	}

	return ""
}
