// Package lznt1 source: https://github.com/Velocidex/go-ntfs/blob/d467c5e7dca0/lznt1.go
package lznt1

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var (
	ErrBlockSize = errors.New("block size invalid")
	ErrDeflate   = errors.New("deflate error")
)

func Detect(b []byte) bool {
	if len(b) < 2 {
		return false
	}

	h := binary.LittleEndian.Uint16(b)

	if h&0x7000 != 0x3000 || int(h&0x0FFF)+5 > len(b) {
		return false // header invalid
	}

	_, err := Deflate(b)
	return err == nil
}

func Deflate(b []byte) ([]byte, error) {
	var buf bytes.Buffer

	for i := 0; ; {
		if len(b) < i+2 {
			break
		}

		chunkOffset := buf.Len()
		blockOffset := i
		blockHeader := binary.LittleEndian.Uint16(b[i:])
		blockSize := int(blockHeader & 0x0FFF)
		blockEnd := blockOffset + blockSize + 3

		// block size invalid
		if blockSize == 0 {
			break
		}

		i += 2

		// block size too small
		if len(b) < i+blockSize {
			return nil, ErrBlockSize
		}

		// block is not compressed
		if blockHeader&0x8000 == 0 {
			buf.Write(b[i : i+blockSize+1])
			i += blockSize + 1
			continue
		}

		// deflate block
		for i < blockEnd {
			header := b[i]
			i++

			for maskIdx := uint8(0); maskIdx < 8 && i < blockEnd; maskIdx++ {
				mask := byte(1 << maskIdx)

				// not masked
				if mask&header == 0 {
					buf.WriteByte(b[i])
					i++
					continue
				}

				v := binary.LittleEndian.Uint16(b[i:])
				d := delta(buf.Len() - chunkOffset - 1)

				symbolOffset := int(v>>(12-d)) + 1
				symbolLength := int(v&(0xFFF>>d)) + 2
				symbolStart := buf.Len() - symbolOffset

				i += 2

				for j := 0; j < symbolLength+1; j++ {
					if buf.Len() <= symbolOffset+j {
						return nil, ErrDeflate
					}

					buf.WriteByte(buf.Bytes()[symbolStart+j])
				}
			}
		}
	}

	return buf.Bytes(), nil
}

func delta(i int) (b byte) {
	for i >= 0x10 {
		i >>= 1
		b += 1
	}
	return
}
