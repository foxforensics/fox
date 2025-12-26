package heap

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

const (
	file  = "text/bible.txt"
	name  = "test"
	size  = 4633983
	count = 31107
	bytes = 1024
	lines = 3
	grep  = 152
)

func TestNew(t *testing.T) {
	h := New(newCtx(), fixture(file))

	if h.Name != name {
		t.Fatal("invalid name")
	}

	if h.Type != types.Regular {
		t.Fatal("invalid type")
	}

	if h.Size != size {
		t.Fatal("invalid size")
	}

	if len(h.MMap()) != size {
		t.Fatal("invalid mmap len")
	}

	if len(h.SMap()) != count {
		t.Fatal("invalid smap len")
	}
}

func TestNewLimitHeadBytes(t *testing.T) {
	ctx := newCtx()
	ctx.Limit.IsHead = true
	ctx.Limit.Bytes = bytes

	h := New(ctx, fixture(file))

	if len(h.MMap()) != bytes {
		t.Fatal("invalid mmap len")
	}
}

func TestNewLimitHeadLines(t *testing.T) {
	ctx := newCtx()
	ctx.Limit.IsHead = true
	ctx.Limit.Lines = lines

	h := New(ctx, fixture(file))

	if len(h.SMap()) != lines {
		t.Fatal("invalid smap len")
	}
}

func TestNewLimitTailBytes(t *testing.T) {
	ctx := newCtx()
	ctx.Limit.IsTail = true
	ctx.Limit.Bytes = bytes

	h := New(ctx, fixture(file))

	if len(h.MMap()) != bytes {
		t.Fatal("invalid mmap len")
	}
}

func TestNewLimitTailLines(t *testing.T) {
	ctx := newCtx()
	ctx.Limit.IsTail = true
	ctx.Limit.Lines = lines

	h := New(ctx, fixture(file))

	if len(h.SMap()) != lines {
		t.Fatal("invalid smap len")
	}
}

func TestNewFilter(t *testing.T) {
	ctx := newCtx()
	ctx.Filter.Regex = regexp.MustCompile("salvation")

	h := New(ctx, fixture(file))

	if len(h.SMap()) != grep {
		t.Fatal("invalid smap len")
	}
}

func TestNewFilterContext(t *testing.T) {
	ctx := newCtx()
	ctx.Filter.Regex = regexp.MustCompile("salvation")
	ctx.Filter.Before = 1
	ctx.Filter.After = 1

	h := New(ctx, fixture(file))

	if len(h.SMap()) != grep*3 {
		t.Fatal("invalid smap len")
	}
}

func TestMMap(t *testing.T) {
	h := New(newCtx(), fixture(file))

	if h.MMap() == nil {
		t.Fatal("mmap nil")
	}
}

func TestSMap(t *testing.T) {
	h := New(newCtx(), fixture(file))

	if h.SMap() == nil {
		t.Fatal("smap nil")
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

	m := h.MMap()

	if m != nil {
		t.Fatal("mmap not nil")
	}

	s := h.SMap()

	if s != nil {
		t.Fatal("smap not nil")
	}
}

func newCtx() *Context {
	return &Context{
		name,
		types.Regular,
		&types.Limits{
			IsHead: false,
			IsTail: false,
			Lines:  0,
			Bytes:  0,
		},
		&types.Filters{
			Regex:  nil,
			Before: 0,
			After:  0,
		},
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
