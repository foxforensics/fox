package ar

import (
	"bytes"
	"io"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/mkrautz/goar"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte{
		0x21, 0x3C, 0x61, 0x72, 0x63, 0x68, 0x3E, 0x0A,
	})
}

func Extract(b []byte, root, _ string) (e []data.Entry) {
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

		e = append(e, data.Entry{
			Path: filepath.Join(root, h.Name),
			Data: buf,
		})
	}

	return
}
