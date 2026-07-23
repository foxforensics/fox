package msi

import (
	"testing"

	"go.foxforensics.eu/fox/v5/internal/pkg/tests"
)

const src = "archives/test.msi"

func BenchmarkDetect(b *testing.B) {
	buf := tests.Fixture(src)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkExtract(b *testing.B) {
	buf := tests.Fixture(src)

	for b.Loop() {
		Extract(buf, "", "")
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

func FuzzExtract(f *testing.F) {
	for _, rnd := range tests.Random() {
		f.Add(rnd)
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic on %x: %v", b, r)
			}
		}()

		_ = Extract(b, "", "")
	})
}

func TestDetect(t *testing.T) {
	if !Detect(tests.Fixture(src)) {
		t.Fatal("not detected")
	}
}

func TestExtract(t *testing.T) {
	e := Extract(tests.Fixture(src), "", "")

	if len(e) != 3 {
		t.Fatal("invalid entry count")
	}

	for _, s := range e {
		if len(s.Path) == 0 {
			t.Fatal("invalid entry path")
		}

		if len(s.Data) == 0 {
			t.Fatal("invalid entry data")
		}
	}
}
