// Package nt source: https://github.com/staaldraad/go-ntlm/blob/master/ntlm/md4/md4.go
package nt

import (
	"hash"

	"golang.org/x/crypto/md4"
	"golang.org/x/text/encoding/unicode"
)

type NT struct {
	md4 hash.Hash
}

func New() hash.Hash {
	h := &NT{md4.New()}
	h.Reset()
	return h
}

func (h *NT) BlockSize() int {
	return md4.BlockSize
}

func (h *NT) Size() int {
	return md4.Size
}

func (h *NT) Reset() {
	h.md4.Reset()
}

func (h *NT) Write(b []byte) (n int, err error) {
	enc := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()

	b, err = enc.Bytes(b)

	if err != nil {
		return 0, err
	}

	return h.md4.Write(b)
}

func (h *NT) Sum(b []byte) []byte {
	return h.md4.Sum(b)
}
