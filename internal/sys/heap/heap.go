package heap

import (
	"bytes"
	"io"
	"sync"

	"go.foxforensics.eu/fox/v5/internal/sys/memory"
)

const block = 4096

type Heap struct {
	sync.RWMutex

	Path   string // heap path
	Part   string // heap part
	Hint   string // heap hint
	Size   uint64 // heap size
	token  uint64
	memory []byte
}

func New(path, part, hint string, token uint64, memory memory.MMap) *Heap {
	return &Heap{
		Path:   path,
		Part:   part,
		Hint:   hint,
		Size:   uint64(len(memory)),
		token:  token,
		memory: memory,
	}
}

func FromPath(path, part string) *Heap {
	return New(path, part, "", 0, nil)
}

func FromData(path string, data []byte) *Heap {
	return New(path, "", "", 0, data)
}

func (h *Heap) String() string {
	return h.Path
}

// Reader must not be retained past the heap scope.
func (h *Heap) Reader() io.ReaderAt {
	h.RLock()
	defer h.RUnlock()

	return bytes.NewReader(h.memory)
}

// Bytes must not be retained past the heap scope.
func (h *Heap) Bytes() []byte {
	h.RLock()
	defer h.RUnlock()

	return h.memory
}

func (h *Heap) IsText() bool {
	h.RLock()
	defer h.RUnlock()

	if h.memory == nil {
		return false
	}

	return !bytes.ContainsRune(h.memory[:min(h.Size, block)], 0)
}

func (h *Heap) Derive(b []byte) {
	h.Lock()
	defer h.Unlock()

	h.Size = uint64(len(b))
	h.token = 0
	h.memory = b
}

func (h *Heap) Free() {
	h.Lock()
	defer h.Unlock()

	h.Size = 0

	if h.token > 0 {
		memory.Free(h.token)
	} else if h.memory != nil {
		h.memory = nil
	}
}
