package nt

import (
	"encoding/binary"
	"hash"
	"unicode/utf16"

	"golang.org/x/crypto/md4"
)

type NT struct {
	md4 hash.Hash
}

func New() hash.Hash {
	h := &NT{md4.New()}
	h.Reset()
	return h
}

func (h *NT) BlockSize() int {
	return md4.BlockSize
}

func (h *NT) Size() int {
	return md4.Size
}

func (h *NT) Reset() {
	h.md4.Reset()
}

func (h *NT) Write(b []byte) (n int, err error) {
	uints := utf16.Encode([]rune(string(b)))
	return len(b), binary.Write(h.md4, binary.LittleEndian, &uints)
}

func (h *NT) Sum(b []byte) []byte {
	return h.md4.Sum(b)
}
