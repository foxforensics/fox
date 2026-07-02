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

	Name   string // heap name
	Size   uint64 // heap size
	Stage  byte
	mapped bool
	memory []byte
}

func New(name string, stage byte, mapped bool, memory memory.MMap) *Heap {
	return &Heap{
		Name:   name,
		Size:   uint64(len(memory)),
		Stage:  stage,
		mapped: mapped,
		memory: memory,
	}
}

func (h *Heap) Clone() *Heap {
	return &Heap{
		Name:   h.Name,
		Size:   h.Size,
		Stage:  h.Stage,
		mapped: false,
		memory: nil,
	}
}

func (h *Heap) String() string {
	return h.Name
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

func (h *Heap) Discard() {
	h.Lock()
	defer h.Unlock()

	h.Size = 0

	if h.mapped {
		memory.Free(h.String())
	} else if h.memory != nil {
		h.memory = h.memory[:0]
	}
}
