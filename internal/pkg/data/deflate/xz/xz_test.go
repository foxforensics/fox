package xz

import (
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

const file = "deflate/fox.xz"

func BenchmarkDetect(b *testing.B) {
	buf := data.FixtureRaw(file)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkDeflate(b *testing.B) {
	buf := data.FixtureRaw(file)

	for b.Loop() {
		_, _ = Deflate(buf)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(data.FixtureRaw(file)) {
		t.Fatal("not detected")
	}
}

func TestDeflate(t *testing.T) {
	buf, err := Deflate(data.FixtureRaw(file))

	if err != nil {
		t.Error(err)
	}

	if !data.Assert(buf) {
		t.Fatal("not deflated")
	}
}
