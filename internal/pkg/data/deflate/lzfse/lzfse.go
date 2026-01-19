package lzfse

import (
	"bytes"
	"io"

	"github.com/aixiansheng/lzfse"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

func Detect(b []byte) bool {
	for _, m := range [][]byte{
		{0x62, 0x76, 0x78, 0x24}, // bvx$
		{0x62, 0x76, 0x78, 0x2D}, // bvx-
		{0x62, 0x76, 0x78, 0x31}, // bvx1
		{0x62, 0x76, 0x78, 0x32}, // bvx2
		{0x62, 0x76, 0x78, 0x6E}, // bvxn
	} {
		if data.HasMagic(b, 0, m) {
			return true
		}
	}

	return false
}

func Deflate(b []byte) ([]byte, error) {
	return io.ReadAll(lzfse.NewReader(bytes.NewReader(b)))
}
