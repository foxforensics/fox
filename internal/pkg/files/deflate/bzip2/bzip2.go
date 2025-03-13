package bzip2

import (
	"bytes"
	"compress/bzip2"
	"io"

	"github.com/cuhsat/fox/v4/internal/pkg/files"
)

func Detect(b []byte) bool {
	return files.HasMagic(b, 0, []byte{
		'B', 'Z', 'h',
	})
}

func Deflate(b []byte) ([]byte, error) {
	return io.ReadAll(bzip2.NewReader(bytes.NewReader(b)))
}
