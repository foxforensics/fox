package reg

import (
	"errors"

	"go.foxforensics.eu/fox/v4/internal/pkg"
)

func Detect(b []byte) bool {
	return pkg.HasMagic(b, 0, []byte{
		'r', 'e', 'g', 'f',
	})
}

func Convert(b []byte) ([]byte, error) {
	return b, errors.New("not implemented")
}
