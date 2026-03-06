package image

import (
	"bytes"
	"hash"
	"image"
	"log"

	"github.com/ajdnik/imghash/v2"
)

const (
	Average Type = iota
	Difference
	Median
	PHash
	WHash
	MarrHildreth
	BlockMean
	PDQ
	RASH
)

type Type int

type Hash struct {
	base imghash.Hasher
	sum  imghash.Hash
}

func New(typ Type) hash.Hash {
	var h imghash.Hasher
	var err error

	switch typ {
	case Average:
		h, err = imghash.NewAverage()
	case Difference:
		h, err = imghash.NewDifference()
	case Median:
		h, err = imghash.NewMedian()
	case PHash:
		h, err = imghash.NewPHash()
	case WHash:
		h, err = imghash.NewWHash()
	case MarrHildreth:
		h, err = imghash.NewMarrHildreth()
	case BlockMean:
		h, err = imghash.NewBlockMean()
	case PDQ:
		h, err = imghash.NewPDQ()
	case RASH:
		h, err = imghash.NewRASH()
	}

	if err != nil {
		log.Fatal(err)
	}

	return &Hash{base: h}
}

func (h *Hash) BlockSize() int {
	return h.sum.Len()
}

func (h *Hash) Size() int {
	return h.sum.Len()
}

func (h *Hash) Reset() {
	// not supported
}

func (h *Hash) Write(b []byte) (n int, err error) {
	img, _, err := image.Decode(bytes.NewReader(b))

	if err != nil {
		return 0, err
	}

	h.sum, err = h.base.Calculate(img)

	if err != nil {
		return 0, err
	}

	return len(b), nil
}

func (h *Hash) Sum(_ []byte) []byte {
	if v, ok := h.sum.(imghash.Binary); ok {
		return v
	}

	return nil
}
