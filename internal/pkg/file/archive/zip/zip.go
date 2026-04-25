package zip

import (
	"bytes"
	"io"
	"log"
	"strings"

	"go.foxforensics.dev/go-zip/pkg/zip"

	"go.foxforensics.dev/fox/v4/internal/pkg/file"
)

func Detect(b []byte) bool {
	for _, m := range [][]byte{
		{'P', 'K', 0x03, 0x04}, // default
		{'P', 'K', 0x03, 0x06}, // empty
		{'P', 'K', 0x03, 0x08}, // spanned
	} {
		if file.HasMagic(b, 0, m) {
			return true
		}
	}

	return false
}

func Extract(b []byte, root, pass string) (e []file.Stream) {
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

		e = append(e, file.Stream{
			Path: file.JoinPart(root, f.Name),
			Data: buf,
		})
	}

	return
}
