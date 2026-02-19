package reg

import (
	"bytes"
	"testing"

	"foxhunt.dev/fox/internal/pkg/data/extract"
	"foxhunt.dev/fox/internal/pkg/test"
)

const file = "ntds/system.zst"

func BenchmarkBootKey(b *testing.B) {
	buf := bytes.NewReader(test.Fixture(file))

	for b.Loop() {
		_, _ = BootKey(buf)
	}
}

func TestBootKey(t *testing.T) {
	buf, err := BootKey(bytes.NewReader(test.Fixture(file)))

	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(buf, extract.BootKey) {
		t.Fatal("invalid bootkey")
	}
}
