package pf

import (
	"encoding/json"
	"testing"

	"foxhunt.dev/fox/internal/pkg/test"
)

const file = "convert/test.pf.zst"

func BenchmarkDetect(b *testing.B) {
	buf := test.Fixture(file)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkConvert(b *testing.B) {
	buf := test.Fixture(file)

	for b.Loop() {
		_, _ = Convert(buf)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(test.Fixture(file)) {
		t.Fatal("not detected")
	}
}

func TestConvert(t *testing.T) {
	buf, err := Convert(test.Fixture(file))

	if err != nil {
		t.Error(err)
	}

	if !json.Valid(buf) {
		t.Fatal("invalid json")
	}
}
