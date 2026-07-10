package br

import (
	"bytes"

	"github.com/andybalholm/brotli"
	"go.foxforensics.eu/fox/v5/library"
)

func Detect(b []byte) bool {
	return library.HasMagic(b, 0, []byte{
		0xCE, 0xB2, 0xCF, 0x81, // Framing Format Signature v3
	})
}

func Deflate(b []byte) ([]byte, error) {
	// remove header
	return library.ReadMax(brotli.NewReader(bytes.NewReader(b[4:])), len(b))
}
