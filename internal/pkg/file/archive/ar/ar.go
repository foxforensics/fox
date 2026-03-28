package ar

import (
	"bytes"
	"io"
	"io/fs"
	"log"

	"github.com/mkrautz/goar"

	"go.foxforensics.dev/fox/v4/internal/pkg/file"
)

func Detect(b []byte) bool {
	return file.HasMagic(b, 0, []byte{
		0x21, 0x3C, 0x61, 0x72, 0x63, 0x68, 0x3E, 0x0A,
	})
}

func Extract(b []byte, root, _ string) (e []file.Stream) {
	r := ar.NewReader(bytes.NewReader(b))

	for {
		h, err := r.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
			break
		}

		if h.Mode&int64(fs.ModeDir) != 0 {
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
