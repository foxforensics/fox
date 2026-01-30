package journal

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

const file = "convert/test.journal.xz"

func BenchmarkDetect(b *testing.B) {
	buf := data.Fixture(file)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkConvert(b *testing.B) {
	buf := data.Fixture(file)

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

	lines := strings.Split(string(buf), "\n")

	if len(lines) != 1923 {
		t.Fatal("invalid length")
	}

	for _, l := range lines {
		if len(l) > 0 && !json.Valid([]byte(l)) {
			t.Fatal("invalid json")
		}
	}
}
