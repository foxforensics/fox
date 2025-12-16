// Package lm source: https://github.com/staaldraad/go-ntlm/blob/master/ntlm/crypto.go
package lm

import (
	"crypto/des"
	"hash"
	"log"
	"strings"
)

const (
	size  = 16
	block = 14
)

type LM struct {
	sum []byte
}

func New() hash.Hash {
	return new(LM)
}

func (h *LM) BlockSize() int {
	return block
}

func (h *LM) Size() int {
	return size
}

func (h *LM) Reset() {
	h.sum = h.sum[:0]
}

func (h *LM) Write(b []byte) (n int, err error) {
	if len(b) > block {
		log.Fatalln("input size to large")
	}

	h.Reset()

	s := strings.ToUpper(string(b))

	if len(s) < block {
		s += strings.Repeat("\x00", block-len(s))
	}

	s1 := desBlock(desKey([]byte(s[:7])))
	s2 := desBlock(desKey([]byte(s[7:])))

	h.sum = append(h.sum, s1...)
	h.sum = append(h.sum, s2...)

	return len(b), nil
}

func (h *LM) Sum(_ []byte) []byte {
	return h.sum
}

func desBlock(k []byte) []byte {
	b := make([]byte, 8)

	c, err := des.NewCipher(k)

	if err != nil {
		log.Fatalln(err)
	}

	c.Encrypt(b, []byte("KGS!@#$%"))

	return b
}

func desKey(b []byte) []byte {
	k := make([]byte, 8)

	k[0] = b[0]
	k[1] = b[0]<<7 | (b[1]&0xff)>>1
	k[2] = b[1]<<6 | (b[2]&0xff)>>2
	k[3] = b[2]<<5 | (b[3]&0xff)>>3
	k[4] = b[3]<<4 | (b[4]&0xff)>>4
	k[5] = b[4]<<3 | (b[5]&0xff)>>5
	k[6] = b[5]<<2 | (b[6]&0xff)>>6
	k[7] = b[6] << 1

	for i := 0; i < len(k); i++ {
		b := k[i]

		if (((b >> 7) ^ (b >> 6) ^ (b >> 5) ^ (b >> 4) ^ (b >> 3) ^ (b >> 2) ^ (b >> 1)) & 0x01) == 0 {
			k[i] = k[i] | byte(0x01)
		} else {
			k[i] = k[i] & byte(0xfe)
		}
	}

	return k
}
