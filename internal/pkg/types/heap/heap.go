package heap

import (
	"bytes"
	"io"
	"math"
	"runtime"
	"sync"

	"github.com/cuhsat/fox/v4/internal/pkg/types/mmap"
)

type Heap struct {
	sync.RWMutex

	Name string // heap name
	Hint string // heap hint
	Time uint64 // heap time
	Size uint64 // heap size

	m mmap.MMap // memory map
}

func New(name, hint string, time, size uint64, m mmap.MMap) *Heap {
	return &Heap{Name: name, Hint: hint, Time: time, Size: size, m: m}
}

func (h *Heap) String() string {
	return h.Name
}

func (h *Heap) Reader() io.ReaderAt {
	h.RLock()
	defer h.RUnlock()
	return bytes.NewReader(h.m)
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
		h.m = nil
	}

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
