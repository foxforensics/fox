package heap

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

const (
	file  = "text/bible.txt"
	name  = "test"
	size  = 4633983
	bytes = 1024
)

func TestNew(t *testing.T) {
	h := New(newCtx(), fixture(file))

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

func TestNewLimitHeadBytes(t *testing.T) {
	ctx := newCtx()
	ctx.Limit.IsHead = true
	ctx.Limit.Bytes = bytes

	h := New(ctx, fixture(file))

	if len(h.Bytes()) != bytes {
		t.Fatal("invalid bytes len")
	}
}

func TestNewLimitTailBytes(t *testing.T) {
	ctx := newCtx()
	ctx.Limit.IsTail = true
	ctx.Limit.Bytes = bytes

	h := New(ctx, fixture(file))

	if len(h.Bytes()) != bytes {
		t.Fatal("invalid bytes len")
	}
}

func TestBytes(t *testing.T) {
	h := New(newCtx(), fixture(file))

	if h.Bytes() == nil {
		t.Fatal("bytes nil")
	}
}

func TestString(t *testing.T) {
	h := New(newCtx(), fixture(file))

	if h.String() != name {
		t.Fatal("string invalid")
	}
}

func TestDiscard(t *testing.T) {
	h := New(newCtx(), fixture(file))
	h.Discard()

	if h.Size > 0 {
		t.Fatal("invalid size")
	}

	m := h.Bytes()

	if m != nil {
		t.Fatal("bytes not nil")
	}
}

func newCtx() *Context {
	return &Context{
		name,
		&types.Limits{},
	}
}

func fixture(name string) []byte {
	const dir = "../../../../testdata"

	_, c, _, ok := runtime.Caller(0)

	if !ok {
		log.Fatalln("runtime error")
	}

	buf, err := os.ReadFile(filepath.Join(filepath.Dir(c), dir, name))

	if err != nil {
		log.Fatalln(err)
	}

	return buf
}
