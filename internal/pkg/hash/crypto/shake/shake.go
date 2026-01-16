package shake

import (
	"crypto/sha3"
)

type Shake struct {
	sha3.SHAKE

	size int
}

func New128() *Shake {
	return &Shake{*sha3.NewSHAKE128(), 16}
}

func New256() *Shake {
	return &Shake{*sha3.NewSHAKE256(), 32}
}

func (s *Shake) Size() int {
	return s.size
}

func (s *Shake) Sum(_ []byte) []byte {
	b := make([]byte, s.size)
	_, _ = s.Read(b)

	return b
}
