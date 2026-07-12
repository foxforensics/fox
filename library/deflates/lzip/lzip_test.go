package lzip

import (
	"testing"

	"go.foxforensics.eu/fox/v5/internal/pkg/tests"
)

const src = "deflates/fox.lz"

func BenchmarkDetect(b *testing.B) {
	buf := tests.Fixture(src)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkDeflate(b *testing.B) {
	buf := tests.Fixture(src)

	for b.Loop() {
		_, _ = Deflate(buf)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(tests.Fixture(src)) {
		t.Fatal("not detected")
	}
}

func TestDeflate(t *testing.T) {
	buf, err := Deflate(tests.Fixture(src))

	if err != nil {
		t.Error(err)
	}

	if !tests.AssertFox(buf) {
		t.Fatal("not deflated")
	}
}
