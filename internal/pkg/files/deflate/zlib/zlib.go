package zlib

import (
	"bytes"
	"io"

	"github.com/klauspost/compress/zlib"

	"github.com/cuhsat/fox/v4/internal/pkg/files"
)

func Detect(b []byte) bool {
	for _, m := range [][]byte{
		{0x78, 0x01}, // no compression
		{0x78, 0x5E}, // fast compression
		{0x78, 0x9C}, // default compression
		{0x78, 0xDA}, // best compression
	} {
		if files.HasMagic(b, 0, m) {
			return true
		}
	}

	return false
}

func Deflate(b []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(b))

	if err != nil {
		return nil, err
	}

	defer func(r io.Closer) {
		_ = r.Close()
	}(r)

	return io.ReadAll(r)
}
