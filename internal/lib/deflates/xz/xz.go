package xz

import (
	"bytes"

	"github.com/ulikunitz/xz"
	"go.foxforensics.eu/fox/v4/internal/lib"
)

func Detect(b []byte) bool {
	return lib.HasMagic(b, 0, []byte{
		0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00,
	})
}

func Deflate(b []byte) ([]byte, error) {
	r, err := xz.NewReader(bytes.NewReader(b))

	if err != nil {
		return b, err
	}

	return lib.ReadMax(r, len(b))
}
