package vhdx

import (
	"encoding/binary"
	"os"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/test"
)

const file = "hunt/test.vhdx.zst"

func BenchmarkDetect(b *testing.B) {
	buf := test.Fixture(file)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkReader(b *testing.B) {
	s := test.FixtureDeflate(file)

	defer func() {
		_ = os.Remove(s)
	}()

	f, err := os.Open(s)

	defer func() {
		_ = f.Close()
	}()

	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		_, _ = Reader(f)
	}

	_ = f.Close()
}

func TestDetect(t *testing.T) {
	if !Detect(test.Fixture(file)) {
		t.Fatal("not detected")
	}
}

func TestReader(t *testing.T) {
	s := test.FixtureDeflate(file)

	defer func() {
		_ = os.Remove(s)
	}()

	f, err := os.Open(s)

	defer func() {
		_ = f.Close()
	}()

	if err != nil {
		t.Fatal(err)
	}

	r, err := Reader(f)

	if err != nil {
		t.Fatal(err)
	}

	b := make([]byte, 2)

	n, err := r.ReadAt(b, 510)

	if err != nil {
		t.Fatal("not read")
	}

	if n != 2 || len(b) != 2 {
		t.Fatal("not read fully")
	}

	if binary.BigEndian.Uint16(b) != 0x55AA {
		t.Fatal("no mbr found")
	}
}
