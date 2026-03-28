package br

import (
	"bytes"
	"io"

	"github.com/andybalholm/brotli"

	"go.foxforensics.dev/fox/v4/internal/pkg/file"
)

func Detect(b []byte) bool {
	return file.HasMagic(b, 0, []byte{
		0xCE, 0xB2, 0xCF, 0x81, // Framing Format Signature v3
	})
}

func Deflate(b []byte) ([]byte, error) {
	return io.ReadAll(brotli.NewReader(bytes.NewReader(b[4:])))
}
