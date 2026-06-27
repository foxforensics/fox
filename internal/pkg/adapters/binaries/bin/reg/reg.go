// Package reg is only implemented for file detection
package reg

import (
	"go.foxforensics.eu/fox/v4/internal/pkg"
)

func Detect(b []byte) bool {
	return pkg.HasMagic(b, 0, []byte{
		'r', 'e', 'g', 'f',
	})
}
