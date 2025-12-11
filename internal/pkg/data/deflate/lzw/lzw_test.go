package lzw

import (
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

const file = "fox.gs.Z"

func BenchmarkDetect(b *testing.B) {
	buf := data.Fixture(file)

	b.ResetTimer()

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkDeflate(b *testing.B) {
	buf := data.Fixture(file)

	b.ResetTimer()

	for b.Loop() {
		_, _ = Deflate(buf)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(data.Fixture(file)) {
		t.Fatal("not detected")
	}
}

func TestDeflate(t *testing.T) {
	buf, err := Deflate(data.Fixture(file))

	if err != nil {
		t.Error(err)
	}

	if !data.Assert(buf) {
		t.Fatal("not deflated")
	}
}
