package xar

import (
	"bytes"
	"io"
	"log"
	"strings"

	"github.com/korylprince/goxar"

	"go.foxforensics.dev/fox/v4/internal/pkg/file"
)

func Detect(b []byte) bool {
	return file.HasMagic(b, 0, []byte{
		'x', 'a', 'r', '!',
	})
}

func Extract(b []byte, root, _ string) (e []file.Stream) {
	r, err := xar.NewReader(nop(bytes.NewReader(b)), int64(len(b)))

	defer func() {
		_ = r.Close()
	}()

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

func extractFile(f *xar.File, root string) (e file.Stream, err error) {
	r, err := f.Open()

	if err != nil {
		return e, err
	}

	defer func(r io.Closer) {
		_ = r.Close()
	}(r)

	e.Path = file.JoinPart(root, f.Name)
	e.Data, err = io.ReadAll(r)

	return e, err
}
