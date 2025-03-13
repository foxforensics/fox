package lzw

import (
	"bytes"
	"compress/lzw"
	"io"

	"github.com/cuhsat/fox/v4/internal/pkg/files"
)

func Detect(b []byte) bool {
	return files.HasMagic(b, 0, []byte{
		0x1F, 0x9D,
	})
}

func Deflate(b []byte) ([]byte, error) {
	r := lzw.NewReader(bytes.NewReader(b), lzw.MSB, 8)
	defer r.Close()

	return io.ReadAll(r)
}
