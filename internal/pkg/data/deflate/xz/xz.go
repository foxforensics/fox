package xz

import (
	"bytes"
	"io"

	"github.com/ulikunitz/xz"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte{
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
