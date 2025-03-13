package br

import (
	"bytes"
	"io"

	"github.com/andybalholm/brotli"

	"github.com/cuhsat/fox/v4/internal/pkg/files"
)

func Detect(b []byte) bool {
	return files.HasMagic(b, 0, []byte{
		0xCE, 0xB2, 0xCF, 0x81,
	})
}

func Deflate(b []byte) ([]byte, error) {
	return io.ReadAll(brotli.NewReader(bytes.NewReader(b)))
}
