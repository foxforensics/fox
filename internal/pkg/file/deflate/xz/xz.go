package xz

import (
	"bytes"
	"io"

	"github.com/ulikunitz/xz"

	"go.foxforensics.dev/fox/v4/internal/pkg/file"
)

func Detect(b []byte) bool {
	return file.HasMagic(b, 0, []byte{
		0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00,
	})
}

func Deflate(b []byte) ([]byte, error) {
	r, err := xz.NewReader(bytes.NewReader(b))

	if err != nil {
		return b, err
	}

	return io.ReadAll(r)
}
