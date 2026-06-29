package zstd

import (
	"bytes"
	"io"

	"github.com/klauspost/compress/zstd"
	"go.foxforensics.eu/fox/v4/internal/pkg/lib"
)

func Detect(b []byte) bool {
	return lib.HasMagic(b, 0, []byte{
		0x28, 0xB5, 0x2F, 0xFD,
	})
}

func Deflate(b []byte) ([]byte, error) {
	r, err := zstd.NewReader(bytes.NewReader(b))

	if err != nil {
		return b, err
	}

	defer r.Close()

	return io.ReadAll(r)
}
