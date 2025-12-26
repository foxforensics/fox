package fson

import (
	"strings"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

const file = "format/fox.fson"

func BenchmarkDetect(b *testing.B) {
	buf := data.Fixture(file)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkFormat(b *testing.B) {
	buf := data.Fixture(file)

	for b.Loop() {
		_ = Format(buf)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(data.Fixture(file)) {
		t.Fatal("not detected")
	}
}

func TestFormat(t *testing.T) {
	buf := Format(data.Fixture(file))

	lines := strings.Split(string(buf), "\n")

	if len(lines) != 4 {
		t.Fatal("invalid length")
	}
}
