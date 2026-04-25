package msi

import (
	"testing"

	"go.foxforensics.dev/fox/v4/internal/pkg/test"
)

const src = "archive/test.msi"

func BenchmarkDetect(b *testing.B) {
	buf := test.Fixture(src)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkExtract(b *testing.B) {
	buf := test.Fixture(src)

	for b.Loop() {
		Extract(buf, "", "")
	}
}

func TestDetect(t *testing.T) {
	if !Detect(test.Fixture(src)) {
		t.Fatal("not detected")
	}
}

func TestExtract(t *testing.T) {
	e := Extract(test.Fixture(src), "", "")

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
