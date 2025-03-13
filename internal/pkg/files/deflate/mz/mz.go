package mz

import (
	"bytes"
	"io"

	"github.com/minio/minlz"

	"github.com/cuhsat/fox/v4/internal/pkg/files"
)

func Detect(b []byte) bool {
	return files.HasMagic(b, 0, []byte{
		0xFF, 0x06, 0x00, 0x00, 0x4D, 0x69, 0x6E, 0x4C, 0x7A,
	})
}

func Deflate(b []byte) ([]byte, error) {
	return io.ReadAll(minlz.NewReader(bytes.NewReader(b)))
}
