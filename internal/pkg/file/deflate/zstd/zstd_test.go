package zstd

import (
	"testing"

	"go.foxforensics.dev/fox/v4/internal/pkg/test"
)

const src = "deflate/fox.zst"

func BenchmarkDetect(b *testing.B) {
	buf := test.FixtureRaw(src)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkDeflate(b *testing.B) {
	buf := test.FixtureRaw(src)

	for b.Loop() {
		_, _ = Deflate(buf)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(test.FixtureRaw(src)) {
		t.Fatal("not detected")
	}
}

func TestDeflate(t *testing.T) {
	buf, err := Deflate(test.FixtureRaw(src))

	if err != nil {
		t.Error(err)
	}

	if !test.Assert(buf) {
		t.Fatal("not deflated")
	}
}
