package bzip2

import (
	"bytes"
	"compress/bzip2"

	"go.foxforensics.eu/fox/v5/library"
)

func Detect(b []byte) bool {
	return library.HasMagic(b, 0, []byte{
		'B', 'Z', 'h',
	})
}

func Deflate(b []byte) ([]byte, error) {
	return library.ReadMax(bzip2.NewReader(bytes.NewReader(b)), len(b))
}
