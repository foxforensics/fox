package snappy

import (
	"bytes"
	"io"

	"github.com/klauspost/compress/snappy"

	"go.foxforensics.dev/fox/v4/internal/pkg/file"
)

func Detect(b []byte) bool {
	return file.HasMagic(b, 0, []byte{
		0xFF, 0x06, 0x00, 0x00, 0x73, 0x4E, 0x61, 0x50, 0x70, 0x59,
	})
}

func Deflate(b []byte) ([]byte, error) {
	return io.ReadAll(snappy.NewReader(bytes.NewReader(b)))
}
