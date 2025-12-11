package tar

import (
	"archive/tar"
	"bytes"
	"io"
	"log"
	"path/filepath"
	"strings"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 257, []byte{
		0x75, 0x73, 0x74, 0x61, 0x72,
	})
}

func Extract(b []byte, root, _ string) (e []data.Entry) {
	r := tar.NewReader(bytes.NewBuffer(b))

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
			Path: filepath.Join(root, h.Name),
			Data: buf,
		})
	}

	return
}
