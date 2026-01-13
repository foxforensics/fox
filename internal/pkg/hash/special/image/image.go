package image

import (
	"bytes"
	"encoding/binary"
	"hash"
	"image"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/corona10/goimagehash"
)

const (
	ahash Kind = iota
	dhash
	phash
)

type Kind int

type Hash struct {
	Base *goimagehash.ImageHash
	Kind Kind
}

func NewAHash() hash.Hash {
	return &Hash{Kind: ahash}
}

func NewDHash() hash.Hash {
	return &Hash{Kind: dhash}
}

func NewPHash() hash.Hash {
	return &Hash{Kind: phash}
}

func (h *Hash) BlockSize() int {
	return h.Base.Bits() / 8
}

func (h *Hash) Size() int {
	return h.Base.Bits() / 8
}

func (h *Hash) Reset() {
	h.Base = nil
}

func (h *Hash) Write(b []byte) (n int, err error) {
	img, _, err := image.Decode(bytes.NewReader(b))

	if err != nil {
		return 0, err
	}

	switch h.Kind {
	case ahash:
		h.Base, err = goimagehash.AverageHash(img)
	case dhash:
		h.Base, err = goimagehash.DifferenceHash(img)
	case phash:
		h.Base, err = goimagehash.PerceptionHash(img)
	}

	if err != nil {
		return 0, err
	}

	return len(b), nil
}

func (h *Hash) Sum(_ []byte) []byte {
	b := make([]byte, h.Size())

	binary.BigEndian.PutUint64(b, h.Base.GetHash())

	return b
}
