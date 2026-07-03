package zlib

import (
	"bytes"
	"log/slog"

	"github.com/klauspost/compress/zlib"
	"go.foxforensics.eu/fox/v4/internal/lib"
)

func Detect(b []byte) bool {
	for _, m := range [][]byte{
		{0x78, 0x01}, // no compression
		{0x78, 0x5E}, // fast compression
		{0x78, 0x9C}, // default compression
		{0x78, 0xDA}, // best compression
	} {
		if lib.HasMagic(b, 0, m) {
			return true
		}
	}

	return false
}

func Deflate(b []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(b))

	if err != nil {
		return b, err
	}

	defer func() {
		if err := r.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()

	return lib.ReadMax(r, len(b))
}
