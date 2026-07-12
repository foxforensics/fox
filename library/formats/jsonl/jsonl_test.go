package jsonl

import (
	"strings"
	"testing"

	"go.foxforensics.eu/fox/v5/internal/pkg/tests"
)

const src = "formats/fox.jsonl"

func BenchmarkDetect(b *testing.B) {
	buf := tests.Fixture(src)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkFormat(b *testing.B) {
	buf := tests.Fixture(src)

	for b.Loop() {
		_, _ = Format(buf)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(tests.Fixture(src)) {
		t.Fatal("not detected")
	}
}

func TestFormat(t *testing.T) {
	buf, err := Format(tests.Fixture(src))

	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(string(buf), "\n")

	if len(lines) != 16 {
		t.Fatal("invalid length")
	}
}
