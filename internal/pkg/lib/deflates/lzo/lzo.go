package lzo

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/rasky/go-lzo"
	"go.foxforensics.eu/fox/v4/internal/pkg/lib"
)

func Detect(b []byte) bool {
	return lib.HasMagic(b, 0, []byte{
		0x89, 0x4C, 0x5A, 0x4F, 0x00, 0x0D, 0x0A, 0x1A, 0x0A,
	})
}

func Deflate(b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	if len(b) < 34 {
		return buf.Bytes(), errors.New("invalid length")
	}

	// remove header
	head := 34 + int(b[33]) + 4

	if len(b) <= head {
		return buf.Bytes(), errors.New("invalid length")
	}

	// remove end
	end := len(b) - 4

	if len(b) <= end {
		return buf.Bytes(), errors.New("invalid length")
	}

	body := b[head:end]

	for {
		ul := int(binary.BigEndian.Uint32(body[0:4]))
		cl := int(binary.BigEndian.Uint32(body[4:8]))

		if len(body) < 12+cl {
			return buf.Bytes(), errors.New("invalid block")
		}

		r := bytes.NewReader(body[12 : 12+cl])

		// decompress every block
		blk, err := lzo.Decompress1X(r, cl, ul)

		if err != nil {
			return buf.Bytes(), err
		}

		buf.Write(blk)

		body = body[12+cl:]

		if len(body) < 8 {
			break
		}
	}

	return buf.Bytes(), nil
}
