package heap

import (
	"bytes"
	"io"
	"sync"

	"go.foxforensics.eu/fox/v4/internal/sys/mmap"
)

type Heap struct {
	sync.RWMutex

	Name string // heap name
	Hint string // heap hint
	Size uint64 // heap size

	m mmap.MMap // memory map
}

func New(name, hint string, m mmap.MMap) *Heap {
	return &Heap{Name: name, Hint: hint, Size: uint64(len(m)), m: m}
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

func (h *Heap) IsText() bool {
	h.RLock()
	defer h.RUnlock()
	b := h.m[:min(h.Size, 4096)]
	return !bytes.ContainsRune(b, 0)
}

func (h *Heap) Discard() {
	defer func() {
		_ = recover()
	}()

	h.Lock()

	// unmap memory
	if h.m != nil {
		mmap.Unmap(h.m)
		h.m = nil
	}

	h.Size = 0
	h.Unlock()

	// runtime.GC()
}
