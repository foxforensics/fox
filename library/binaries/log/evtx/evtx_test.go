package evtx

import (
	"encoding/json"
	"strings"
	"testing"

	"go.foxforensics.eu/fox/v5/internal/pkg/tests"
)

const src = "binaries/test.evtx"

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

	lines := strings.Split(string(buf), "\n")

	if len(lines) == 0 {
		t.Fatal("invalid length")
	}

	for _, l := range lines {
		if len(l) > 0 && !json.Valid([]byte(l)) {
			t.Fatal("invalid json")
		}
	}
}
