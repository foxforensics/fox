package szip

import (
	"bytes"
	"io"
	"log"

	"github.com/bodgit/sevenzip"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte{
		0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C,
	})
}

func Extract(b []byte, root, pass string) (e []data.Entry) {
	r, err := sevenzip.NewReaderWithPassword(bytes.NewReader(b), int64(len(b)), pass)

	if err != nil {
		log.Println(err)
		return
	}

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
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

// https://github.com/bodgit/sevenzip?tab=readme-ov-file#why-is-my-code-running-so-slow
func extractFile(f *sevenzip.File, root string) (e data.Entry, err error) {
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
