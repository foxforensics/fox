// Package pe source: https://github.com/saferwall/pe/blob/v1.5.7/helper.go#L552
package pe

import (
	"encoding/binary"
	"hash"
)

const (
	size  = 4
	block = 4
)

type PE struct {
	sum uint32
}

func New() hash.Hash {
	return new(PE)
}

func (h *PE) BlockSize() int {
	return block
}

func (h *PE) Size() int {
	return size
}

func (h *PE) Reset() {
	h.sum = 0
}

func (h *PE) Write(b []byte) (n int, err error) {
	var sum uint64
	var blk uint32

	r := uint32(len(b)) % 4
	l := uint32(len(b))

	if r > 0 {
		l += 4 - r
		b = append(b, make([]byte, 4-r)...)
	}

	for i := uint64(0); i < uint64(l); i += 4 {
		blk = binary.LittleEndian.Uint32(b[i:])
		sum = (sum & 0xffffffff) + uint64(blk) + (sum >> 32)

		if sum > 0x100000000 {
			sum = (sum & 0xffffffff) + (sum >> 32)
		}
	}

	sum = (sum & 0xffff) + (sum >> 16)
	sum = sum + (sum >> 16)
	sum = sum & 0xffff
	sum += uint64(l)

	h.sum = uint32(sum)

	return len(b), nil
}

func (h *PE) Sum(_ []byte) []byte {
	b := make([]byte, size)

	binary.LittleEndian.PutUint32(b, h.sum)

	return b
}
