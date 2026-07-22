package xml

import (
	"strings"
	"testing"

	"go.foxforensics.eu/fox/v5/internal/pkg/tests"
)

const src = "formats/fox.xml"

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

func FuzzDetect(f *testing.F) {
	for _, rnd := range tests.Random() {
		f.Add(rnd)
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic on %x: %v", b, r)
			}
		}()

		_ = Detect(b)
	})
}

func FuzzFormat(f *testing.F) {
	for _, rnd := range tests.Random() {
		f.Add(rnd)
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic on %x: %v", b, r)
			}
		}()

		_, _ = Format(b)
	})
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

	if len(lines) != 6 {
		t.Fatal("invalid length")
	}
}
