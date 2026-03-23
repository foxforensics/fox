package lzo

import (
	"bytes"
	"encoding/binary"

	"github.com/rasky/go-lzo"

	"github.com/cuhsat/fox/v4/internal/pkg/file"
)

func Detect(b []byte) bool {
	return file.HasMagic(b, 0, []byte{
		0x89, 0x4C, 0x5A, 0x4F, 0x00, 0x0D, 0x0A, 0x1A, 0x0A,
	})
}

func Deflate(b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	// remove header
	// 89 4c 5a 4f 00 0d 0a 1a 0a
	// 10 40
	// 20 a0
	// 09 40
	// 02
	// 01
	// 03 00 00 01
	// 00 00 81 a4
	// 69 5a c7 48
	// 00 00 00 00
	// 07
	// 66 6f 78 2e 74 78 74
	// 61 e8 07 3a
	head := 34 + int(b[33]) + 4

	// remove end
	// 00 00 00 00
	end := len(b) - 4

	body := b[head:end]

	for {
		ul := int(binary.BigEndian.Uint32(body[0:4]))
		cl := int(binary.BigEndian.Uint32(body[4:8]))

		r := bytes.NewReader(body[12 : 12+cl])

		// decompress every block
		blk, err := lzo.Decompress1X(r, cl, ul)

		if err != nil {
			return buf.Bytes(), err
		}

		buf.Write(blk)

		body = body[12+cl:]

		if len(body) == 0 {
			break
		}
	}

	return buf.Bytes(), nil
}
