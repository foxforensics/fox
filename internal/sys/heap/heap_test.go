package heap

import (
	"testing"

	"go.foxforensics.eu/fox/v4/internal/test"
)

const (
	file = "texts/bible.txt"
	name = "test"
	size = 4633983
)

func TestNew(t *testing.T) {
	h := FromData(name, test.Fixture(file))

	defer h.DeAlloc()

	if h.Size != size {
		t.Fatal("invalid size")
	}

	if h.String() != name {
		t.Fatal("invalid name")
	}

	if len(h.Bytes()) != size {
		t.Fatal("invalid bytes len")
	}
}

func TestIsText(t *testing.T) {
	h := FromData(name, test.Fixture(file))

	defer h.DeAlloc()

	if !h.IsText() {
		t.Fatal("not text")
	}
}

func TestBytes(t *testing.T) {
	h := FromData(name, test.Fixture(file))

	defer h.DeAlloc()

	if h.Bytes() == nil {
		t.Fatal("bytes nil")
	}
}

func TestString(t *testing.T) {
	h := FromData(name, test.Fixture(file))

	defer h.DeAlloc()

	if h.String() != name {
		t.Fatal("string invalid")
	}
}

func TestDiscard(t *testing.T) {
	h := FromData(name, test.Fixture(file))
	h.DeAlloc()

	if h.Size > 0 {
		t.Fatal("invalid size")
	}

	m := h.Bytes()

	if m != nil {
		t.Fatal("bytes not nil")
	}
}
