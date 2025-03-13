package blake3

import (
	"github.com/zeebo/blake3"
)

type Blake3 struct {
	blake3.Hasher

	size int
}

func New256() *Blake3 {
	return &Blake3{*blake3.New(), 32}
}

func New512() *Blake3 {
	return &Blake3{*blake3.New(), 64}
}

func (b *Blake3) Sum(_ []byte) []byte {
	s := make([]byte, b.size)
	d := b.Hasher.Digest()

	_, _ = d.Read(s)

	return s
}
