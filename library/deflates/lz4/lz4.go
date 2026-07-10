package lz4

import (
	"bytes"

	"github.com/pierrec/lz4/v4"
	"go.foxforensics.eu/fox/v5/library"
)

func Detect(b []byte) bool {
	return library.HasMagic(b, 0, []byte{
		0x04, 0x22, 0x4D, 0x18,
	})
}

func Deflate(b []byte) ([]byte, error) {
	return library.ReadMax(lz4.NewReader(bytes.NewReader(b)), len(b))
}
