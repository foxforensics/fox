package reg

import (
	"bytes"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/file/extract"
	"github.com/cuhsat/fox/v4/internal/pkg/test"
)

const src = "dump/test.reg.zst"

func BenchmarkBootKey(b *testing.B) {
	buf := bytes.NewReader(test.Fixture(src))

	for b.Loop() {
		_, _ = BootKey(buf)
	}
}

func TestBootKey(t *testing.T) {
	buf, err := BootKey(bytes.NewReader(test.Fixture(src)))

	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(buf, extract.BootKey) {
		t.Fatal("invalid bootkey")
	}
}
