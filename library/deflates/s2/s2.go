package s2

import (
	"bytes"

	"github.com/klauspost/compress/s2"
	"go.foxforensics.eu/fox/v4/library"
)

func Detect(b []byte) bool {
	return library.HasMagic(b, 0, []byte{
		0xFF, 0x06, 0x00, 0x00, 0x53, 0x32, 0x73, 0x54, 0x77, 0x4F,
	})
}

func Deflate(b []byte) ([]byte, error) {
	return library.ReadMax(s2.NewReader(bytes.NewReader(b)), len(b))
}
