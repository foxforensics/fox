package lzip

import (
	"bytes"
	"io"

	"github.com/sorairolake/lzip-go"

	"github.com/cuhsat/fox/v4/internal/pkg/files"
)

func Detect(b []byte) bool {
	return files.HasMagic(b, 0, []byte{
		0x4C, 0x5A, 0x49, 0x50,
	})
}

func Deflate(b []byte) ([]byte, error) {
	r, err := lzip.NewReader(bytes.NewReader(b))

	if err != nil {
		return nil, err
	}

	return io.ReadAll(r)
}
