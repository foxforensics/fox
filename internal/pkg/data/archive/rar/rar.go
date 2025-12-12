package rar

import (
	"bytes"
	"io"
	"log"
	"strings"

	"github.com/nwaples/rardecode"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte{
		0x52, 0x61, 0x72, 0x21, 0x1A, 0x07,
	})
}

func Extract(b []byte, root, pass string) (e []data.Entry) {
	r, err := rardecode.NewReader(bytes.NewBuffer(b), pass)

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

		e = append(e, data.Entry{
			Path: data.AddStream(root, h.Name),
			Data: buf,
		})
	}

	return
}
