package heap

import (
	"errors"
	"log"
	"math"
	"runtime"
	"sync"
	"syscall"

	"github.com/cuhsat/go-mmap"

	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

type Heap struct {
	sync.RWMutex
	Name string    // heap name
	Hint string    // heap hint
	Size uint64    // heap size
	mmap mmap.MMap // memory map
}

func New(name, hint string, m mmap.MMap, l *types.Limits) *Heap {
	return &Heap{
		Name: name,
		Hint: hint,
		Size: uint64(len(m)),
		mmap: l.Reduce(m),
	}
}

func (h *Heap) String() string {
	return h.Name
}

func (h *Heap) Bytes() []byte {
	h.RLock()
	defer h.RUnlock()
	return h.mmap
}

func (h *Heap) Discard() {
	h.Lock()

	// try to unmap original area
	err := h.mmap.Unmap()

	if err != nil && !errors.Is(err, syscall.EINVAL) {
		log.Println(err)
	}

	h.Size = 0
	h.mmap = nil

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
