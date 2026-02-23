// Package imphash based on https://github.com/omarghader/pefile-go/blob/master/pe/pe.go
package imphash

import (
	"crypto/md5"
	"hash"
	"strings"

	"github.com/cuhsat/fox/v4/internal/pkg/hash/fuzzy"
)

type ImpHash struct {
	sort bool
	buf  []string
}

func New() hash.Hash {
	return &ImpHash{sort: false}
}

func NewStable() hash.Hash {
	return &ImpHash{sort: true}
}

func (h *ImpHash) BlockSize() int {
	return md5.BlockSize // from underlying MD5
}

func (h *ImpHash) Size() int {
	return md5.Size
}

func (h *ImpHash) Reset() {
	h.buf = h.buf[:0]
}

func (h *ImpHash) Write(b []byte) (n int, err error) {
	h.buf, err = fuzzy.GetImports(b, h.sort)

	return len(b), err

}

func (h *ImpHash) Sum(_ []byte) []byte {
	sum := md5.Sum([]byte(strings.Join(h.buf, ",")))

	return sum[:]
}
