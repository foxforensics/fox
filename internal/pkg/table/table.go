package table

import (
	"bytes"
	"sync"

	"github.com/sourcegraph/conc/iter"
	"go.foxforensics.dev/hasher/hash"
	"go.foxforensics.dev/wordlist"
)

var hashesLm sync.Map
var hashesNt sync.Map

func Build(b []byte) (err error) {
	// use built-in wordlist for rainbow tables
	if b == nil {
		b, err = wordlist.Deflate()

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
