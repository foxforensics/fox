package s2

import (
	"bytes"
	"io"

	"github.com/klauspost/compress/s2"

	"github.com/cuhsat/fox/v4/internal/pkg/files"
)

func Detect(b []byte) bool {
	return files.HasMagic(b, 0, []byte{
		0xFF, 0x06, 0x00, 0x00, 0x53, 0x32, 0x73, 0x54, 0x77, 0x4F,
	})
}

func Deflate(b []byte) ([]byte, error) {
	return io.ReadAll(s2.NewReader(bytes.NewReader(b)))
}
