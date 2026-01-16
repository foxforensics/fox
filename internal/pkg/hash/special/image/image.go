package image

import (
	"bytes"
	"encoding/binary"
	"hash"
	"image"
	"log"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/corona10/goimagehash"
)

type Hash struct {
	base *goimagehash.ImageHash
	kind goimagehash.Kind
}

func NewAHash() hash.Hash {
	return &Hash{kind: goimagehash.AHash}
}

func NewDHash() hash.Hash {
	return &Hash{kind: goimagehash.DHash}
}

func NewPHash() hash.Hash {
	return &Hash{kind: goimagehash.PHash}
}

func (h *Hash) BlockSize() int {
	return h.base.Bits() / 8
}

func (h *Hash) Size() int {
	return h.base.Bits() / 8
}

func (h *Hash) Reset() {
	h.base = goimagehash.NewImageHash(0, h.kind)
}

func (h *Hash) Write(b []byte) (n int, err error) {
	img, _, err := image.Decode(bytes.NewReader(b))

	if err != nil {
		return 0, err
	}

	switch h.kind {
	case goimagehash.AHash:
		h.base, err = goimagehash.AverageHash(img)
	case goimagehash.DHash:
		h.base, err = goimagehash.DifferenceHash(img)
	case goimagehash.PHash:
		h.base, err = goimagehash.PerceptionHash(img)
	default:
		log.Fatalln("unknown kind")
	}

	if err != nil {
		return 0, err
	}

	return len(b), nil
}

func (h *Hash) Sum(_ []byte) []byte {
	b := make([]byte, h.Size())

	binary.BigEndian.PutUint64(b, h.base.GetHash())

	return b
}
