// Package library provides support for different file types.
//
// Source:
// - https://github.com/frizb/FirmwareReverseEngineering/blob/master/IdentifyingCompressionAlgorithms.md
// - https://en.wikipedia.org/wiki/List_of_file_signatures
package library

import (
	"bytes"
	"io"

	"go.foxforensics.eu/fox/v5/internal/pkg"
)

type Chunk struct {
	Path string // Chunk path
	Data []byte // Chunk data
}

type Detect func([]byte) bool

type Format func([]byte) ([]byte, error)

type Convert func([]byte) ([]byte, error)

type Deflate func([]byte) ([]byte, error)

type Extract func([]byte, string, string) []Chunk

func HasMagic(b []byte, off int, m []byte) bool {
	if len(b) < off+len(m) {
		return false
	}

	return bytes.Equal(b[off:off+len(m)], m)
}

func ReadMax(r io.Reader, size int) ([]byte, error) {
	// prevent zip bombs, that deflate exorbitant amounts of data
	return io.ReadAll(io.LimitReader(r, int64(size*pkg.MaxFactor)))
}
