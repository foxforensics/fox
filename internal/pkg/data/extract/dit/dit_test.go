package dit

import (
	"bytes"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/data/extract"
	"github.com/cuhsat/fox/v4/internal/pkg/test"
)

const file = "ntds/ntds.dit.zst"
const dump = "ntds/dump.golden"

func BenchmarkExtract(b *testing.B) {
	buf := test.Fixture(file)

	for b.Loop() {
		_, _, _ = Extract(buf, extract.BootKey)
	}
}

func TestExtract(t *testing.T) {
	var buf bytes.Buffer

	rec, key, err := Extract(test.Fixture(file), extract.BootKey)

	if err != nil {
		t.Error(err)
	}

	if len(rec) == 0 {
		t.Fatal("no records")
	}

	if len(key) == 0 {
		t.Fatal("ne keys")
	}

	for _, r := range rec {
		buf.WriteString(r.String())
		buf.WriteByte('\n')
	}

	if !bytes.Equal(buf.Bytes(), test.Fixture(dump)) {
		t.Fatal("wrong hashes")
	}
}
