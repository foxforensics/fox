package lz4

import (
	"bytes"
	"io"

	"github.com/pierrec/lz4/v4"

	"go.foxforensics.dev/fox/v4/internal/pkg/file"
)

func Detect(b []byte) bool {
	return file.HasMagic(b, 0, []byte{
		0x04, 0x22, 0x4D, 0x18,
	})
}

func Deflate(b []byte) ([]byte, error) {
	return io.ReadAll(lz4.NewReader(bytes.NewReader(b)))
}
