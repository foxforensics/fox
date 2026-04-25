package pf

import (
	"bytes"
	"encoding/json"

	"www.velocidex.com/golang/go-prefetch"

	"go.foxforensics.dev/fox/v4/internal/pkg/file"
)

func Detect(b []byte) bool {
	for _, v := range []struct {
		off int
		buf []byte
	}{
		{off: 4, buf: []byte{'S', 'C', 'C', 'A'}},  // uncompressed
		{off: 0, buf: []byte{'M', 'A', 'M', 0x04}}, // LZX compressed
	} {
		if file.HasMagic(b, v.off, v.buf) {
			return true
		}
	}

	return false
}

func Convert(b []byte) ([]byte, error) {
	pi, err := prefetch.LoadPrefetch(bytes.NewReader(b))

	if err != nil {
		return b, err
	}

	return json.MarshalIndent(pi, "", "  ")
}
