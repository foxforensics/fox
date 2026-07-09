package minlz

import (
	"bytes"

	"github.com/minio/minlz"
	"go.foxforensics.eu/fox/v4/library"
)

func Detect(b []byte) bool {
	return library.HasMagic(b, 0, []byte{
		0xFF, 0x06, 0x00, 0x00, 0x4D, 0x69, 0x6E, 0x4C, 0x7A,
	})
}

func Deflate(b []byte) ([]byte, error) {
	return library.ReadMax(minlz.NewReader(bytes.NewReader(b)), len(b))
}
