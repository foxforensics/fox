package zstd

import (
	"bytes"
	"io"

	"github.com/klauspost/compress/zstd"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

func Detect(b []byte) bool {
	for _, m := range [][]byte{
		{0x1E, 0xB5, 0x2F, 0xFD}, // v0.1
		{0x22, 0xB5, 0x2F, 0xFD}, // v0.2
		{0x23, 0xB5, 0x2F, 0xFD}, // v0.3
		{0x24, 0xB5, 0x2F, 0xFD}, // v0.4
		{0x25, 0xB5, 0x2F, 0xFD}, // v0.5
		{0x26, 0xB5, 0x2F, 0xFD}, // v0.6
		{0x27, 0xB5, 0x2F, 0xFD}, // v0.7
		{0x28, 0xB5, 0x2F, 0xFD}, // v0.8
	} {
		if data.HasMagic(b, 0, m) {
			return true
		}
	}

	return false
}

func Deflate(b []byte) ([]byte, error) {
	r, err := zstd.NewReader(bytes.NewReader(b))

	if err != nil {
		return b, err
	}

	defer r.Close()

	return io.ReadAll(r)
}
