// Package table source: https://github.com/danielmiessler/SecLists/blob/master/Passwords/Common-Credentials/
package table

import (
	"bytes"
	"io"
	"sync"

	_ "embed"

	"github.com/sourcegraph/conc/iter"
	"github.com/ulikunitz/xz"
	"go.foxforensics.dev/hasher/hash"
)

//go:embed wordlist.xz
var wordlist []byte
var hashesLm sync.Map
var hashesNt sync.Map

func Build(b []byte) error {
	if b == nil {
		r, err := xz.NewReader(bytes.NewReader(wordlist))

		if err != nil {
			return err
		}

		b, err = io.ReadAll(r)

		if err != nil {
			return err
		}
	}

	iter.ForEach(bytes.Split(b, []byte{'\n'}), func(b *[]byte) {
		hashesNt.Store(hash.MustSum(hash.NT, *b), string(*b))
		hashesLm.Store(hash.MustSum(hash.LM, *b), string(*b))
	})

	return nil
}

func Lookup(sum string) string {
	if v, ok := hashesNt.Load(sum); ok {
		return v.(string)
	}

	if v, ok := hashesLm.Load(sum); ok {
		return v.(string)
	}

	return ""
}
