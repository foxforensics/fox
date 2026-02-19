package lzip

import (
	"bytes"
	"io"

	"github.com/sorairolake/lzip-go"

	"foxhunt.dev/fox/internal/pkg/data"
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte{
		0x4C, 0x5A, 0x49, 0x50,
	})
}

func Deflate(b []byte) ([]byte, error) {
	r, err := lzip.NewReader(bytes.NewReader(b))

	if err != nil {
		return b, err
	}

	return io.ReadAll(r)
}
