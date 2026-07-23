package lzfse

import (
	"testing"

	"go.foxforensics.eu/fox/v5/internal/pkg/tests"
)

const src = "deflates/fox.lzfse"

func BenchmarkDetect(b *testing.B) {
	buf := tests.Fixture(src)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkDeflate(b *testing.B) {
	buf := tests.Fixture(src)

	for b.Loop() {
		_, _ = Deflate(buf)
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

func FuzzDeflate(f *testing.F) {
	for _, rnd := range tests.Random() {
		f.Add(rnd)
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic on %x: %v", b, r)
			}
		}()

		_, _ = Deflate(b)
	})
}

func TestDetect(t *testing.T) {
	if !Detect(tests.Fixture(src)) {
		t.Fatal("not detected")
	}
}

func TestDeflate(t *testing.T) {
	buf, err := Deflate(tests.Fixture(src))

	if err != nil {
		t.Error(err)
	}

	if !tests.AssertFox(buf) {
		t.Fatal("not deflated")
	}
}
