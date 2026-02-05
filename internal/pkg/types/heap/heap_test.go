package heap

import (
	"os"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/test"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

const (
	file = "text/bible.txt"
	name = "test"
	hint = "text"
	size = 4633983
	page = 1024
)

func TestFromData(t *testing.T) {
	h := FromData(name, hint, size, test.Fixture(file), new(types.Limits))

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

func TestFromFile(t *testing.T) {
	f, err := os.Open(test.FixtureFile(file))

	if err != nil {
		t.Fatal(err)
	}

	h := FromFile(name, hint, size, f, f)

	defer h.Discard()

	if h.Name != name {
		t.Fatal("invalid name")
	}

	if h.Size != size {
		t.Fatal("invalid size")
	}

	if len(h.Bytes()) != 0 {
		t.Fatal("invalid bytes len")
	}
}

func TestNewLimitHeadBytes(t *testing.T) {
	h := FromData(name, hint, size, test.Fixture(file), &types.Limits{
		IsHead: true,
		Bytes:  page,
	})

	defer h.Discard()

	if len(h.Bytes()) != page {
		t.Fatal("invalid bytes len")
	}
}

func TestNewLimitTailBytes(t *testing.T) {
	h := FromData(name, hint, size, test.Fixture(file), &types.Limits{
		IsTail: true,
		Bytes:  page,
	})

	defer h.Discard()

	if len(h.Bytes()) != page {
		t.Fatal("invalid bytes len")
	}
}

func TestBytes(t *testing.T) {
	h := FromData(name, hint, size, test.Fixture(file), new(types.Limits))

	defer h.Discard()

	if h.Bytes() == nil {
		t.Fatal("bytes nil")
	}
}

func TestString(t *testing.T) {
	h := FromData(name, hint, size, test.Fixture(file), new(types.Limits))

	defer h.Discard()

	if h.String() != name {
		t.Fatal("string invalid")
	}
}

func TestDiscard(t *testing.T) {
	h := FromData(name, hint, size, test.Fixture(file), new(types.Limits))
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
	h := FromData(name, hint, size, test.Fixture(file), new(types.Limits))

	defer h.Discard()

	if Entropy(h.Bytes()) != 0.5758916753281705 {
		t.Fatal("entropy wrong")
	}
}
