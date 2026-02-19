package jsonl

import (
	"strings"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/test"
)

const file = "format/fox.jsonl"

func BenchmarkDetect(b *testing.B) {
	buf := test.Fixture(file)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkFormat(b *testing.B) {
	buf := test.Fixture(file)

	for b.Loop() {
		_, _ = Format(buf)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(test.Fixture(file)) {
		t.Fatal("not detected")
	}
}

func TestFormat(t *testing.T) {
	buf, err := Format(test.Fixture(file))

	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(string(buf), "\n")

	if len(lines) != 16 {
		t.Fatal("invalid length")
	}
}
