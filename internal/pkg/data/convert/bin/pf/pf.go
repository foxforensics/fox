package pf

import (
	"bytes"
	"encoding/json"

	"www.velocidex.com/golang/go-prefetch"

	"foxhunt.dev/fox/internal/pkg/data"
)

func Detect(b []byte) bool {
	for _, m := range [][]byte{
		{'S', 'C', 'C', 'A'},  // uncompressed
		{'M', 'A', 'M', 0x04}, // LZX compressed
	} {
		if data.HasMagic(b, 0, m) {
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
