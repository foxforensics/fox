// Package file source:
// https://github.com/frizb/FirmwareReverseEngineering/blob/master/IdentifyingCompressionAlgorithms.md
// https://en.wikipedia.org/wiki/List_of_file_signatures
package file

import (
	"bytes"
	"math"
)

type Stream struct {
	Path string // Stream path
	Data []byte // Stream data
}

type Detect func([]byte) bool

type Format func([]byte) ([]byte, error)

type Convert func([]byte) ([]byte, error)

type Deflate func([]byte) ([]byte, error)

type Extract func([]byte, string, string) []Stream

func HasMagic(b []byte, off int, m []byte) bool {
	if len(b) < off+len(m) {
		return false
	}

	return bytes.Equal(b[off:off+len(m)], m)
}

func Entropy(b []byte) float64 {
	var v [256]float64
	var e float64

	for _, b := range b {
		v[b]++
	}

	l := float64(len(b))

	for i := range 256 {
		if v[i] != 0 {
			f := v[i] / l
			e -= f * math.Log2(f)
		}
	}

	return e
}
