// Package reg is only implemented for file detection
package reg

import (
	"go.foxforensics.eu/fox/v4/library"
)

func Detect(b []byte) bool {
	return library.HasMagic(b, 0, []byte{
		'r', 'e', 'g', 'f',
	})
}
