package mft

import (
	"encoding/json"
	"strings"
	"testing"

	"go.foxforensics.eu/fox/v5/internal/pkg/tests"
)

const src = "binaries/test.mft"

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

func FuzzConvert(f *testing.F) {
	for _, rnd := range tests.Random() {
		f.Add(rnd)
	}

	f.Fuzz(func(t *testing.T, b []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic on %x: %v", b, r)
			}
		}()

		_, _ = Convert(b)
	})
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
