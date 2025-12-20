// Package data source:
// https://github.com/frizb/FirmwareReverseEngineering/blob/master/IdentifyingCompressionAlgorithms.md
// https://en.wikipedia.org/wiki/List_of_file_signatures
package data

import "bytes"

type Entry struct {
	Path string // Entry path
	Data []byte // Entry data
}

type Format func([]byte, int) ([]byte, error)

type Deflate func([]byte) ([]byte, error)

type Extract func([]byte, string, string) []Entry

func HasMagic(b []byte, off int, m []byte) bool {
	if len(b) < off+len(m) {
		return false
	}

	return bytes.Equal(b[off:off+len(m)], m)
}
