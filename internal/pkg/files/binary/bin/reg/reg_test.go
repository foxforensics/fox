package reg

import (
	"testing"

	"go.foxforensics.eu/fox/v4/internal/test"
)

const src = "binary/test.dat"

func BenchmarkDetect(b *testing.B) {
	buf := test.Fixture(src)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(test.Fixture(src)) {
		t.Fatal("not detected")
	}
}
