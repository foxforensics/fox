package snappy

import (
	"bytes"

	"github.com/klauspost/compress/snappy"
	"go.foxforensics.eu/fox/v5/library"
)

func Detect(b []byte) bool {
	return library.HasMagic(b, 0, []byte{
		0xFF, 0x06, 0x00, 0x00, 0x73, 0x4E, 0x61, 0x50, 0x70, 0x59,
	})
}

func Deflate(b []byte) ([]byte, error) {
	return library.ReadMax(snappy.NewReader(bytes.NewReader(b)), len(b))
}
