package cab

import (
	"bytes"
	"io"
	"log"
	"path/filepath"
	"strings"

	"github.com/google/go-cabfile/cabfile"

	"github.com/cuhsat/fox/v4/internal/pkg/files"
)

func Detect(b []byte) bool {
	return files.HasMagic(b, 0, []byte{
		0x4D, 0x53, 0x43, 0x46,
	})
}

func Extract(b []byte, root, _ string) (e []files.Entry) {
	r, err := cabfile.New(bytes.NewReader(b))

	if err != nil {
		log.Println(err)
		return
	}

	for _, s := range r.FileList() {
		if strings.HasSuffix(s, "/") {
			continue
		}

		h, err := r.Content(s)

		if err != nil {
			log.Println(err)
			continue
		}

		buf, err := io.ReadAll(h)

		if err != nil {
			log.Println(err)
			continue
		}

		e = append(e, files.Entry{
			Path: filepath.Join(root, s),
			Data: buf,
		})
	}

	return
}
