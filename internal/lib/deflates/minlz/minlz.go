package minlz

import (
	"bytes"
	"io"

	"github.com/minio/minlz"
	"go.foxforensics.eu/fox/v4/internal/lib"
)

func Detect(b []byte) bool {
	return lib.HasMagic(b, 0, []byte{
		0xFF, 0x06, 0x00, 0x00, 0x4D, 0x69, 0x6E, 0x4C, 0x7A,
	})
}

func Deflate(b []byte) ([]byte, error) {
	return io.ReadAll(minlz.NewReader(bytes.NewReader(b)))
}
