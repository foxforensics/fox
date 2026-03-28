package dit

import (
	"bytes"
	"testing"

	"go.foxforensics.dev/fox/v4/internal/pkg/file/extract"
	"go.foxforensics.dev/fox/v4/internal/pkg/test"
)

const src = "dump/test.dit.zst"
const dst = "dump/dump.golden"

func BenchmarkExtract(b *testing.B) {
	buf := test.Fixture(src)

	for b.Loop() {
		_, _, _ = Extract(buf, extract.BootKey)
	}
}

func TestExtract(t *testing.T) {
	var buf bytes.Buffer

	rec, key, err := Extract(test.Fixture(src), extract.BootKey)

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

	if !bytes.Equal(buf.Bytes(), test.Fixture(dst)) {
		t.Fatal("wrong hashes")
	}
}
