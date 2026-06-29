package bzip2

import (
	"bytes"
	"compress/bzip2"
	"io"

	"go.foxforensics.eu/fox/v4/internal/pkg/lib"
)

func Detect(b []byte) bool {
	return lib.HasMagic(b, 0, []byte{
		'B', 'Z', 'h',
	})
}

func Deflate(b []byte) ([]byte, error) {
	return io.ReadAll(bzip2.NewReader(bytes.NewReader(b)))
}
