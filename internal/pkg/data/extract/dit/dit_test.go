package dit

import (
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/data/extract"
	"github.com/cuhsat/fox/v4/internal/pkg/test"
)

const file = "ntds/ntds.dit.zst"

func BenchmarkExtract(b *testing.B) {
	buf := test.Fixture(file)

	for b.Loop() {
		_, _ = Extract(buf, extract.BootKey)
	}
}

func TestExtract(t *testing.T) {
	rec, err := Extract(test.Fixture(file), extract.BootKey)

	if err != nil {
		t.Error(err)
	}

	if len(rec) == 0 {
		t.Fatal("no records")
	}

	for _, r := range rec {
		println(r.String())
	}
}
