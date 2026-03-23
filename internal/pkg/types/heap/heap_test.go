package heap

import (
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/test"
)

const (
	file = "text/bible.txt"
	name = "test"
	hint = "text"
	size = 4633983
)

func TestNew(t *testing.T) {
	h := New(name, hint, size, test.Fixture(file))

	defer h.Discard()

	if h.Name != name {
		t.Fatal("invalid name")
	}

	if h.Size != size {
		t.Fatal("invalid size")
	}

	if len(h.Bytes()) != size {
		t.Fatal("invalid bytes len")
	}
}

func TestBytes(t *testing.T) {
	h := New(name, hint, size, test.Fixture(file))

	defer h.Discard()

	if h.Bytes() == nil {
		t.Fatal("bytes nil")
	}
}

func TestString(t *testing.T) {
	h := New(name, hint, size, test.Fixture(file))

	defer h.Discard()

	if h.String() != name {
		t.Fatal("string invalid")
	}
}

func TestDiscard(t *testing.T) {
	h := New(name, hint, size, test.Fixture(file))
	h.Discard()

	if h.Size > 0 {
		t.Fatal("invalid size")
	}

	m := h.Bytes()

	if m != nil {
		t.Fatal("bytes not nil")
	}
}

func TestEntropy(t *testing.T) {
	h := New(name, hint, size, test.Fixture(file))

	defer h.Discard()

	if Entropy(h.Bytes()) != 4.607133402625364 {
		t.Fatal("entropy wrong")
	}
}
