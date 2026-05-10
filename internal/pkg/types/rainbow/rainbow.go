// Package rainbow source: https://github.com/danielmiessler/SecLists/blob/master/Passwords/Common-Credentials/probable-v2_top-12000.txt
package rainbow

import (
	"bytes"
	"io"
	"strings"
	"sync"

	_ "embed"

	"github.com/klauspost/compress/zstd"
	"github.com/sourcegraph/conc/iter"
	"go.foxforensics.dev/hasher/hash"
)

//go:embed rainbow.zst
var wordlist []byte

var table sync.Map

func Build(parallel int) error {
	r, err := zstd.NewReader(bytes.NewReader(wordlist))

	if err != nil {
		return err
	}

	defer r.Close()

	b, err := io.ReadAll(r)

	if err != nil {
		return err
	}

	it := iter.Iterator[string]{
		MaxGoroutines: parallel,
	}

	it.ForEach(strings.Split(string(b), "\n"), func(s *string) {
		table.Store(hash.MustSum(hash.NT, []byte(*s)), *s)
	})

	return nil
}

func Lookup(sum string) string {
	if v, ok := table.Load(sum); ok {
		return v.(string)
	}

	return ""
}
