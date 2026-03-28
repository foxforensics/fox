package xml

import (
	"strings"
	"testing"

	"go.foxforensics.dev/fox/v4/internal/pkg/test"
)

const src = "format/fox.xml"

func BenchmarkDetect(b *testing.B) {
	buf := test.Fixture(src)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkFormat(b *testing.B) {
	buf := test.Fixture(src)

	for b.Loop() {
		_, _ = Format(buf)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(test.Fixture(src)) {
		t.Fatal("not detected")
	}
}

func TestFormat(t *testing.T) {
	buf, err := Format(test.Fixture(src))

	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(string(buf), "\n")

	if len(lines) != 6 {
		t.Fatal("invalid length")
	}
}
