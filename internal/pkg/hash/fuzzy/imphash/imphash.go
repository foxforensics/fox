// Package imphash based on https://github.com/omarghader/pefile-go/blob/master/pe/pe.go
package imphash

import (
	"bytes"
	"crypto/md5"
	"debug/pe"
	"errors"
	"fmt"
	"hash"
	"strings"

	intern "github.com/cuhsat/fox/v4/internal/pkg/data/convert/bin/pe"
)

var ErrNotSupported = errors.New("file type not supported")

type ImpHash struct {
	v []string
}

func New() hash.Hash {
	return new(ImpHash)
}

func (h *ImpHash) BlockSize() int {
	return md5.BlockSize // from underlying MD5
}

func (h *ImpHash) Size() int {
	return md5.Size
}

func (h *ImpHash) Reset() {
	h.v = h.v[:0]
}

func (h *ImpHash) Write(b []byte) (n int, err error) {
	if !intern.Detect(b) {
		return 0, ErrNotSupported
	}

	f, err := pe.NewFile(bytes.NewReader(b))

	if err != nil {
		return 0, err
	}

	defer func(f *pe.File) {
		_ = f.Close()
	}(f)

	iat, err := f.ImportedSymbols()

	if err != nil {
		return 0, err
	}

	rep := strings.NewReplacer(".dll", "", ".ocx", "", ".sys", "")

	for _, e := range iat {
		if !strings.Contains(e, ":") {
			continue
		}

		p := strings.Split(rep.Replace(strings.ToLower(e)), ":")

		h.v = append(h.v, fmt.Sprintf("%s.%s", p[1], p[0]))
	}

	return len(b), nil
}

func (h *ImpHash) Sum(_ []byte) []byte {
	sum := md5.Sum([]byte(strings.Join(h.v, ",")))

	return sum[:]
}
