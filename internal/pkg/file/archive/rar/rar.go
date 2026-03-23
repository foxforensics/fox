package rar

import (
	"bytes"
	"io"
	"log"
	"strings"

	"github.com/nwaples/rardecode/v2"

	"github.com/cuhsat/fox/v4/internal/pkg/file"
)

func Detect(b []byte) bool {
	return file.HasMagic(b, 0, []byte{
		0x52, 0x61, 0x72, 0x21, 0x1A, 0x07,
	})
}

func Extract(b []byte, root, pass string) (e []file.Stream) {
	r, err := rardecode.NewReader(bytes.NewBuffer(b), rardecode.Password(pass))

	if err != nil {
		log.Println(err)
		return
	}

	for {
		h, err := r.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
			break
		}

		if strings.HasSuffix(h.Name, "/") {
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
