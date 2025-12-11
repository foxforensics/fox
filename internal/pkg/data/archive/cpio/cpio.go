package cpio

import (
	"bytes"
	"io"
	"log"
	"path/filepath"

	"github.com/cavaliergopher/cpio"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

func Detect(b []byte) bool {
	for _, m := range [][]byte{
		{0x30, 0x37, 0x30, 0x37, 0x30, 0x31}, // SRV4
		{0x30, 0x37, 0x30, 0x37, 0x30, 0x32}, // SRV4 with CRC
	} {
		if data.HasMagic(b, 0, m) {
			return true
		}
	}

	return false
}

func Extract(b []byte, root, _ string) (e []data.Entry) {
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

		e = append(e, data.Entry{
			Path: filepath.Join(root, h.Name),
			Data: buf,
		})
	}

	return
}
