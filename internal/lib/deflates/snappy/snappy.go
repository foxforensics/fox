package snappy

import (
	"bytes"

	"github.com/klauspost/compress/snappy"
	"go.foxforensics.eu/fox/v4/internal/lib"
)

func Detect(b []byte) bool {
	return lib.HasMagic(b, 0, []byte{
		0xFF, 0x06, 0x00, 0x00, 0x73, 0x4E, 0x61, 0x50, 0x70, 0x59,
	})
}

func Deflate(b []byte) ([]byte, error) {
	return lib.ReadMax(snappy.NewReader(bytes.NewReader(b)), len(b))
}
