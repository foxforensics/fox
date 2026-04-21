package djb2

import (
	"encoding/binary"
	"hash"
)

const (
	size  = 8
	block = 8
	start = 5381
)

type Djb2 struct {
	sum uint64
}

func New() hash.Hash {
	return &Djb2{sum: start}
}

func (h *Djb2) BlockSize() int {
	return block
}

func (h *Djb2) Size() int {
	return size
}

func (h *Djb2) Reset() {
	h.sum = start
}

func (h *Djb2) Write(b []byte) (n int, err error) {
	for i := 0; i < len(b); i++ {
		h.sum = ((h.sum << 5) + h.sum) + uint64(b[i])
	}

	return len(b), nil
}

func (h *Djb2) Sum(_ []byte) []byte {
	b := make([]byte, size)

	binary.LittleEndian.PutUint64(b, h.sum)

	return b
}
