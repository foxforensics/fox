package zstd

import (
	"testing"

	"foxhunt.dev/fox/internal/pkg/test"
)

const file = "deflate/fox.zst"

func BenchmarkDetect(b *testing.B) {
	buf := test.FixtureRaw(file)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkDeflate(b *testing.B) {
	buf := test.FixtureRaw(file)

	for b.Loop() {
		_, _ = Deflate(buf)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(test.FixtureRaw(file)) {
		t.Fatal("not detected")
	}
}

func TestDeflate(t *testing.T) {
	buf, err := Deflate(test.FixtureRaw(file))

	if err != nil {
		t.Error(err)
	}

	if !test.Assert(buf) {
		t.Fatal("not deflated")
	}
}
