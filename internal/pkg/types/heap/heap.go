package heap

import (
	"bytes"
	"io"
	"math"
	"runtime"
	"sync"

	"foxhunt.dev/fox/internal/pkg/types"
	"foxhunt.dev/fox/internal/pkg/types/mmap"
)

type Heap struct {
	sync.RWMutex

	Name string // heap name
	Hint string // heap hint
	Size uint64 // heap size

	m mmap.MMap   // memory map
	r io.ReaderAt // file reader
	c []io.Closer // file closer
}

func FromData(name, hint string, size uint64, m mmap.MMap, l *types.Limits) *Heap {
	return &Heap{Name: name, Hint: hint, Size: size, m: l.Reduce(m)}
}

func FromFile(name, hint string, size uint64, r io.ReaderAt, c ...io.Closer) *Heap {
	return &Heap{Name: name, Hint: hint, Size: size, c: c, r: r}
}

func (h *Heap) String() string {
	return h.Name
}

func (h *Heap) Reader() io.ReaderAt {
	h.RLock()
	defer h.RUnlock()

	if len(h.c) == 0 {
		return bytes.NewReader(h.m)
	}

	return h.r
}

func (h *Heap) Bytes() []byte {
	h.RLock()
	defer h.RUnlock()
	return h.m
}

func (h *Heap) Discard() {
	h.Lock()

	// unmap memory
	if h.m != nil {
		mmap.Unmap(h.m)
	}

	// close files
	for _, f := range h.c {
		_ = f.Close()
	}

	h.c = h.c[:0]
	h.r = nil
	h.m = nil

	h.Size = 0

	h.Unlock()

	runtime.GC()
}

func Entropy(block []byte) float64 {
	var a [256]float64
	var v float64

	for _, b := range block {
		a[b]++
	}

	l := float64(len(block))

	for i := range 256 {
		if a[i] != 0 {
			f := a[i] / l
			v -= f * math.Log2(f)
		}
	}

	v /= 8

	return v
}
