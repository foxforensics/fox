package sdhash

import (
	"strings"

	"github.com/eciavatta/sdhash"
)

type SDHash struct {
	f sdhash.SdbfFactory
	s sdhash.Sdbf
}

func New() *SDHash {
	return new(SDHash)
}

func (sd *SDHash) Reset() {
	sd.f = nil
	sd.s = nil
}

func (sd *SDHash) BlockSize() int {
	return sdhash.BlockSize
}

func (sd *SDHash) Size() int {
	return int(sd.s.Size())
}

func (sd *SDHash) Sum(_ []byte) []byte {
	sd.s = sd.f.Compute()

	return []byte(strings.TrimRight(sd.s.String(), "\n"))
}

func (sd *SDHash) Write(b []byte) (int, error) {
	var err error

	sd.f, err = sdhash.CreateSdbfFromBytes(b)

	return len(b), err
}
