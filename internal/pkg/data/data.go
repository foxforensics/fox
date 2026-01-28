// Package data source:
// https://github.com/frizb/FirmwareReverseEngineering/blob/master/IdentifyingCompressionAlgorithms.md
// https://en.wikipedia.org/wiki/List_of_file_signatures
package data

import "bytes"

type Stream struct {
	Path string // Stream path
	Data []byte // Stream data
}

type Detect func([]byte) bool

type Format func([]byte) ([]byte, error)

type Ingest func([]byte) ([]byte, error)

type Convert func([]byte) ([]byte, error)

type Deflate func([]byte) ([]byte, error)

type Extract func([]byte, string, string) []Stream

func HasMagic(b []byte, off int, m []byte) bool {
	if len(b) < off+len(m) {
		return false
	}

	return bytes.Equal(b[off:off+len(m)], m)
}
