package ese

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/test"
)

const ese = "convert/test.ese.zst"
const dit = "ntds/ntds.dit.zst"

func BenchmarkDetect(b *testing.B) {
	buf := test.Fixture(ese)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkConvert(b *testing.B) {
	buf := test.Fixture(ese)

	for b.Loop() {
		_, _ = Convert(buf)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(test.Fixture(ese)) {
		t.Fatal("not detected")
	}
}

func TestConvert(t *testing.T) {
	buf, err := Convert(test.Fixture(ese))

	if err != nil {
		t.Error(err)
	}

	if !json.Valid(buf) {
		t.Fatal("invalid format")
	}
}

func TestExtract(t *testing.T) {
	bootKey, _ := hex.DecodeString("13d20976d63ea5e836036ec8bc68d6eb")

	buf, err := Extract(test.Fixture(dit), bootKey)

	if err != nil {
		t.Error(err)
	}

	println(string(buf))

	if len(buf) == 0 {
		t.Fatal("invalid result")
	}
}
