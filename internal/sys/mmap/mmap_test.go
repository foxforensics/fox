package mmap

import (
	"os"
	"testing"

	"go.foxforensics.eu/fox/v4/internal/test"
)

func BenchmarkMMap(b *testing.B) {
	v := test.FixtureFile("texts/bible.txt")

	f, err := os.Open(v)

	if err != nil {
		b.Fatal(err)
	}

	defer func() {
		_ = f.Close()
	}()

	for b.Loop() {
		_, _ = Map(f)
	}
}

func TestMMap(t *testing.T) {
	v := test.FixtureFile("texts/bible.txt")

	f, err := os.Open(v)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = f.Close()
	}()

	m, err := Map(f)

	if err != nil {
		t.Fatal(err)
	}

	defer Unmap(m)

	if len(m) != 4633983 {
		t.Fatal("wrong size")
	}

	Unmap(m)
}
