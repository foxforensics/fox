package pe

import (
	"encoding/json"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

const file = "convert/fox.exe"

func BenchmarkDetect(b *testing.B) {
	buf := data.Fixture(file)

	b.ResetTimer()

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkConvert(b *testing.B) {
	buf := data.Fixture(file)

	b.ResetTimer()

	for b.Loop() {
		_, _ = Convert(buf)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(data.Fixture(file)) {
		t.Fatal("not detected")
	}
}

func TestConvert(t *testing.T) {
	buf, err := Convert(data.Fixture(file))

	if err != nil {
		t.Error(err)
	}

	if !json.Valid(buf) {
		t.Fatal("invalid json")
	}
}
