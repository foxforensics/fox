package xar

import (
	"bytes"
	"io"
	"log"
	"strings"

	"github.com/korylprince/goxar"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte{
		'x', 'a', 'r', '!',
	})
}

func Extract(b []byte, root, _ string) (e []data.Entry) {
	r, err := xar.NewReader(nop(bytes.NewReader(b)), int64(len(b)))

	if err != nil {
		log.Println(err)
		return
	}

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "/") {
			continue
		}

		if v, err := extractFile(f, root); err == nil {
			e = append(e, v)
		} else {
			log.Println(err)
		}
	}

	return
}

func extractFile(f *xar.File, root string) (e data.Entry, err error) {
	r, err := f.Open()

	if err != nil {
		return e, err
	}

	defer func(r io.Closer) {
		_ = r.Close()
	}(r)

	e.Path = data.JoinPart(root, f.Name)
	e.Data, err = io.ReadAll(r)

	return e, err
}
