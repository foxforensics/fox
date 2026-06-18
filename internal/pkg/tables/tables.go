package tables

import (
	"bytes"
	"io"
	"maps"
	"sync"

	"github.com/sourcegraph/conc/iter"
	"go.foxforensics.eu/hasher/hash"
	"go.foxforensics.eu/wordlist"
)

var Threads = 2 // default

var tables = make(map[string]*table, 2)

type table struct {
	m sync.Map
}

func newTable(algo string, w [][]byte) *table {
	t := new(table)

	iter.Iterator[[]byte]{
		MaxGoroutines: Threads,
	}.ForEach(w, func(b *[]byte) {
		t.m.Store(hash.MustSum(algo, *b), string(*b))
	})

	return t
}

func (t *table) Lookup(s string) string {
	if v, ok := t.m.Load(s); ok {
		return v.(string)
	}

	return ""
}

func Build(b []byte, algos ...string) (int, error) {
	// use built-in wordlist for rainbow tables
	if b == nil {
		r, err := wordlist.Reader()

		if err != nil {
			return 0, err
		}

		b, err = io.ReadAll(r)

		if err != nil {
			return 0, err
		}
	}

	w := bytes.Split(b, []byte{'\n'})

	for _, algo := range algos {
		tables[algo] = newTable(algo, w)
	}

	return len(w), nil
}

func Lookup(s string) (string, string) {
	for k, t := range maps.All(tables) {
		if v := t.Lookup(s); len(v) > 0 {
			return k, v
		}
	}

	return "", ""
}
