package evtx

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/test"
)

const src = "convert/test.evtx.zst"

func BenchmarkDetect(b *testing.B) {
	buf := test.Fixture(src)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkConvert(b *testing.B) {
	buf := test.Fixture(src)

	for b.Loop() {
		_, _ = Convert(buf)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(test.Fixture(src)) {
		t.Fatal("not detected")
	}
}

func TestConvert(t *testing.T) {
	buf, err := Convert(test.Fixture(src))

	if err != nil {
		t.Error(err)
	}

	lines := strings.Split(string(buf), "\n")

	if len(lines) != 920 {
		t.Fatal("invalid length")
	}

	for _, l := range lines {
		if len(l) > 0 && !json.Valid([]byte(l)) {
			t.Fatal("invalid json")
		}
	}
}
