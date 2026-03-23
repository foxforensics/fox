package cpio

import (
	"bytes"
	"io"
	"log"

	"github.com/cavaliergopher/cpio"

	"github.com/cuhsat/fox/v4/internal/pkg/file"
)

func Detect(b []byte) bool {
	for _, m := range [][]byte{
		{0x30, 0x37, 0x30, 0x37, 0x30, 0x31}, // SRV4
		{0x30, 0x37, 0x30, 0x37, 0x30, 0x32}, // SRV4 with CRC
	} {
		if file.HasMagic(b, 0, m) {
			return true
		}
	}

	return false
}

func Extract(b []byte, root, _ string) (e []file.Stream) {
	r := cpio.NewReader(bytes.NewBuffer(b))

	for {
		h, err := r.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
			break
		}

		if !h.Mode.IsRegular() {
			continue
		}

		buf, err := io.ReadAll(r)

		if err != nil {
			log.Println(err)
			continue
		}

		e = append(e, file.Stream{
			Path: file.JoinPart(root, h.Name),
			Data: buf,
		})
	}

	return
}
