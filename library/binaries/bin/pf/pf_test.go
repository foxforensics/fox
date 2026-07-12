package pf

import (
	"encoding/json"
	"testing"

	"go.foxforensics.eu/fox/v5/internal/pkg/tests"
)

const src = "binaries/test.pf"

func BenchmarkDetect(b *testing.B) {
	buf := tests.Fixture(src)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkConvert(b *testing.B) {
	buf := tests.Fixture(src)

	for b.Loop() {
		_, _ = Convert(buf)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(tests.Fixture(src)) {
		t.Fatal("not detected")
	}
}

func TestConvert(t *testing.T) {
	buf, err := Convert(tests.Fixture(src))

	if err != nil {
		t.Error(err)
	}

	if !json.Valid(buf) {
		t.Fatal("invalid json")
	}
}
