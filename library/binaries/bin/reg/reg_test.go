package reg

import (
	"testing"

	"go.foxforensics.eu/fox/v5/internal/pkg/tests"
)

const src = "binaries/test.reg"

func BenchmarkDetect(b *testing.B) {
	buf := tests.Fixture(src)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(tests.Fixture(src)) {
		t.Fatal("not detected")
	}
}
