package heap

import (
	"bytes"
	"io"
	"sync"

	"go.foxforensics.eu/fox/v4/internal/sys/memory"
)

const block = 4096

type Heap struct {
	sync.RWMutex

	Path   string // heap path
	Part   string // heap part
	Size   uint64 // heap size
	Stage  byte
	mapped bool
	memory []byte
}

func New(path, part string, stage byte, mapped bool, memory memory.MMap) *Heap {
	return &Heap{
		Path:   path,
		Part:   part,
		Size:   uint64(len(memory)),
		Stage:  stage,
		mapped: mapped,
		memory: memory,
	}
}

func FromPath(path, part string) *Heap {
	return New(path, part, 0, false, nil)
}

func FromData(path string, data []byte) *Heap {
	return New(path, "", 0, false, data)
}

func (h *Heap) String() string {
	return h.Path
}

func (h *Heap) Reader() io.ReaderAt {
	h.RLock()
	defer h.RUnlock()

	return bytes.NewReader(h.memory)
}

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

func (h *Heap) ReAlloc(b []byte) {
	h.Lock()
	defer h.Unlock()

	h.Size = uint64(len(b))
	h.mapped = false
	h.memory = b
}

func (h *Heap) DeAlloc() {
	h.Lock()
	defer h.Unlock()

	h.Size = 0

	if h.mapped {
		memory.Free(h.String())
	} else if h.memory != nil {
		h.memory = nil
	}
}
