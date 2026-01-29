package gzip

import (
	"bytes"
	"io"

	"github.com/klauspost/compress/gzip"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte{
		0x1F, 0x8B, 0x08,
	})
}

func Deflate(b []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(b))

	if err != nil {
		return b, err
	}

	defer func(r *gzip.Reader) {
		_ = r.Close()
	}(r)

	return io.ReadAll(r)
}
