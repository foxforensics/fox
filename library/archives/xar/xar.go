package xar

import (
	"bytes"
	"log/slog"
	"strings"

	"github.com/korylprince/goxar"
	"go.foxforensics.eu/fox/v5/internal/pkg"
	"go.foxforensics.eu/fox/v5/library"
)

func Detect(b []byte) bool {
	return library.HasMagic(b, 0, []byte{
		'x', 'a', 'r', '!',
	})
}

func Extract(b []byte, root, _ string) (e []library.Chunk) {
	r, err := xar.NewReader(nop(bytes.NewReader(b)), int64(len(b)))

	if err != nil {
		slog.Error(err.Error())
		return
	}

	defer func() {
		if err := r.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "/") {
			continue
		}

		if v, err := extractFile(f, root, len(b)); err == nil {
			e = append(e, v)
		} else {
			slog.Error(err.Error())
		}
	}

	return
}

func extractFile(f *xar.File, root string, size int) (e library.Chunk, err error) {
	r, err := f.Open()

	if err != nil {
		return e, err
	}

	defer func() {
		if err := r.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()

	e.Path = pkg.JoinPart(root, f.Name)
	e.Data, err = library.ReadMax(r, size)

	return e, err
}
