package mmap

import (
	"os"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/test"
)

func BenchmarkMMap(b *testing.B) {
	v := test.FixtureFile("text/bible.txt")

	f, err := os.Open(v)

	if err != nil {
		b.Fatal(err)
	}

	defer func() {
		_ = f.Close()
	}()

	for b.Loop() {
		Map(f)
	}
}

func TestMMap(t *testing.T) {
	v := test.FixtureFile("text/bible.txt")

	f, err := os.Open(v)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = f.Close()
	}()

	m := Map(f)

	defer Unmap(m)

	if len(m) != 4633983 {
		t.Fatal("wrong size")
	}

	Unmap(m)
}
