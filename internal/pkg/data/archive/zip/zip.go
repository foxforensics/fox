package zip

import (
	"bytes"
	"io"
	"log"
	"strings"

	"github.com/cuhsat/zip/pkg/zip"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

func Detect(b []byte) bool {
	for _, m := range [][]byte{
		{'P', 'K', 0x03, 0x04}, // default
		{'P', 'K', 0x03, 0x06}, // empty
		{'P', 'K', 0x03, 0x08}, // spanned
	} {
		if data.HasMagic(b, 0, m) {
			return true
		}
	}

	return false
}

func Extract(b []byte, root, pass string) (e []data.Entry) {
	r, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))

	if err != nil {
		log.Println(err)
		return
	}

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "/") {
			continue
		}

		if len(pass) > 0 {
			f.SetPassword(pass)
		}

		a, err := f.Open()

		if err != nil {
			log.Println(err)
			continue
		}

		buf, err := io.ReadAll(a)

		_ = a.Close()

		if err != nil {
			log.Println(err)
			continue
		}

		e = append(e, data.Entry{
			Path: data.JoinPart(root, f.Name),
			Data: buf,
		})
	}

	return
}
