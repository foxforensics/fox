package lzfse

import (
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/test"
)

const file = "deflate/fox.lzfse"

func BenchmarkDetect(b *testing.B) {
	buf := test.Fixture(file)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkDeflate(b *testing.B) {
	buf := test.Fixture(file)

	for b.Loop() {
		_, _ = Deflate(buf)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(test.Fixture(file)) {
		t.Fatal("not detected")
	}
}

func TestDeflate(t *testing.T) {
	buf, err := Deflate(test.Fixture(file))

	if err != nil {
		t.Error(err)
	}

	if !test.Assert(buf) {
		t.Fatal("not deflated")
	}
}
