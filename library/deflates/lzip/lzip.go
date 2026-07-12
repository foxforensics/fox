package lzip

import (
	"bytes"

	"github.com/sorairolake/lzip-go"
	"go.foxforensics.eu/fox/v5/library"
)

func Detect(b []byte) bool {
	return library.HasMagic(b, 0, []byte{
		0x4C, 0x5A, 0x49, 0x50,
	})
}

func Deflate(b []byte) ([]byte, error) {
	r, err := lzip.NewReader(bytes.NewReader(b))

	if err != nil {
		return b, err
	}

	return library.ReadMax(r, len(b))
}
