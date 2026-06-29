// Package reg is only implemented for file detection
package reg

import (
	"go.foxforensics.eu/fox/v4/internal/pkg/lib"
)

func Detect(b []byte) bool {
	return lib.HasMagic(b, 0, []byte{
		'r', 'e', 'g', 'f',
	})
}
