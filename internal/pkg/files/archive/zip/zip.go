package zip

import (
	"bytes"
	"io"
	"log"
	"path/filepath"
	"strings"

	"github.com/cuhsat/zip/pkg/zip"

	"github.com/cuhsat/fox/v4/internal/pkg/files"
)

func Detect(b []byte) bool {
	return files.HasMagic(b, 0, []byte{
		0x50, 0x4B, 0x03, 0x04,
	})
}

func Extract(b []byte, root, pass string) (e []files.Entry) {
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

		e = append(e, files.Entry{
			Path: filepath.Join(root, f.Name),
			Data: buf,
		})
	}

	return
}
