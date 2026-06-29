package gzip

import (
	"bytes"
	"io"
	"log/slog"

	"github.com/klauspost/compress/gzip"
	"go.foxforensics.eu/fox/v4/internal/pkg/lib"
)

func Detect(b []byte) bool {
	return lib.HasMagic(b, 0, []byte{
		0x1F, 0x8B, 0x08,
	})
}

func Deflate(b []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(b))

	if err != nil {
		return b, err
	}

	defer func() {
		if err := r.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()

	return io.ReadAll(r)
}
