package pe

import (
	"encoding/json"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

const file = "binary/fox.exe"

func BenchmarkDetect(b *testing.B) {
	buf := data.Fixture(file)

	b.ResetTimer()

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkFormat(b *testing.B) {
	buf := data.Fixture(file)

	b.ResetTimer()

	for b.Loop() {
		_, _ = Format(buf, 0)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(data.Fixture(file)) {
		t.Fatal("not detected")
	}
}

func TestFormat(t *testing.T) {
	buf, err := Format(data.Fixture(file), 0)

	if err != nil {
		t.Error(err)
	}

	if !json.Valid(buf) {
		t.Fatal("invalid json")
	}
}
