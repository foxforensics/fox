package lz4

import (
	"testing"

	"go.foxforensics.dev/fox/v4/internal/pkg/test"
)

const src = "deflate/fox.lz4"

func BenchmarkDetect(b *testing.B) {
	buf := test.Fixture(src)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkDeflate(b *testing.B) {
	buf := test.Fixture(src)

	for b.Loop() {
		_, _ = Deflate(buf)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(test.Fixture(src)) {
		t.Fatal("not detected")
	}
}

func TestDeflate(t *testing.T) {
	buf, err := Deflate(test.Fixture(src))

	if err != nil {
		t.Error(err)
	}

	if !test.Assert(buf) {
		t.Fatal("not deflated")
	}
}
