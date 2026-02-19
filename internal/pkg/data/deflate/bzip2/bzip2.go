package bzip2

import (
	"bytes"
	"compress/bzip2"
	"io"

	"foxhunt.dev/fox/internal/pkg/data"
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte{
		'B', 'Z', 'h',
	})
}

func Deflate(b []byte) ([]byte, error) {
	return io.ReadAll(bzip2.NewReader(bytes.NewReader(b)))
}
