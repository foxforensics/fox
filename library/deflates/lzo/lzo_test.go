package lzo

import (
	"testing"

	"go.foxforensics.eu/fox/v5/internal/test"
)

const src = "deflates/fox.lzo"

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

	if !test.AssertFox(buf) {
		t.Fatal("not deflated")
	}
}
