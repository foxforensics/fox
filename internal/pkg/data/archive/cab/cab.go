package cab

import (
	"bytes"
	"io"
	"log"
	"strings"

	"github.com/secDre4mer/go-cab"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte{
		0x4D, 0x53, 0x43, 0x46,
	})
}

func Extract(b []byte, root, _ string) (e []data.Entry) {
	r, err := cab.Open(bytes.NewReader(b), int64(len(b)))

	if err != nil {
		log.Println(err)
		return
	}

	for _, f := range r.Files {
		if strings.HasSuffix(f.Name, "/") {
			continue
		}

		h, err := f.Open()

		if err != nil {
			log.Println(err)
			continue
		}

		buf, err := io.ReadAll(h)

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
