package memory

import (
	"os"
	"testing"

	"go.foxforensics.eu/fox/v5/internal/pkg/tests"
)

func BenchmarkMMap(b *testing.B) {
	v := tests.FixtureFile("texts/bible.txt")

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
	v := tests.FixtureFile("texts/bible.txt")

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

	defer func() {
		if err = Unmap(m); err != nil {
			t.Fatal(err)
		}
	}()

	if len(m) != 4633983 {
		t.Fatal("wrong size")
	}
}
