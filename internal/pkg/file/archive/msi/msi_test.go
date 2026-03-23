package msi

import (
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/test"
)

const file = "archive/test.msi"

func BenchmarkDetect(b *testing.B) {
	buf := test.Fixture(file)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkExtract(b *testing.B) {
	buf := test.Fixture(file)

	for b.Loop() {
		Extract(buf, "", "")
	}
}

func TestDetect(t *testing.T) {
	if !Detect(test.Fixture(file)) {
		t.Fatal("not detected")
	}
}

func TestExtract(t *testing.T) {
	e := Extract(test.Fixture(file), "", "")

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
