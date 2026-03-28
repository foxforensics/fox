package msi

import (
	"bytes"
	"io"
	"log"

	"go.foxforensics.dev/go-msi/pkg/msi"

	"github.com/cuhsat/fox/v4/internal/pkg/file"
)

func Detect(b []byte) bool {
	return file.HasMagic(b, 0, []byte{
		0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1,
	})
}

func Extract(b []byte, root, _ string) (e []file.Stream) {
	r, err := msi.Open(bytes.NewReader(b))

	if err != nil {
		log.Println(err)
		return
	}

	streams := r.Streams()

	for {
		name := streams.Next()

		if len(name) == 0 {
			break
		}

		str, err := r.ReadStream(name)

		if err != nil {
			log.Println(err)
			continue
		}

		buf, err := io.ReadAll(str)

		if err != nil {
			log.Println(err)
			continue
		}

		e = append(e, file.Stream{
			Path: file.JoinPart(root, name),
			Data: buf,
		})
	}

	return
}
