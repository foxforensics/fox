package lzfse

import (
	"bytes"

	"github.com/aixiansheng/lzfse"
	"go.foxforensics.eu/fox/v5/library"
)

func Detect(b []byte) bool {
	for _, m := range [][]byte{
		{0x62, 0x76, 0x78, 0x24}, // bvx$
		{0x62, 0x76, 0x78, 0x2D}, // bvx-
		{0x62, 0x76, 0x78, 0x31}, // bvx1
		{0x62, 0x76, 0x78, 0x32}, // bvx2
		{0x62, 0x76, 0x78, 0x6E}, // bvxn
	} {
		if library.HasMagic(b, 0, m) {
			return true
		}
	}

	return false
}

func Deflate(b []byte) ([]byte, error) {
	return library.ReadMax(lzfse.NewReader(bytes.NewReader(b)), len(b))
}
