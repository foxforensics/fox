package ar

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"log/slog"

	"github.com/mkrautz/goar"
	"go.foxforensics.eu/fox/v5/internal/sys"
	"go.foxforensics.eu/fox/v5/library"
)

func Detect(b []byte) bool {
	return library.HasMagic(b, 0, []byte{
		0x21, 0x3C, 0x61, 0x72, 0x63, 0x68, 0x3E, 0x0A,
	})
}

func Extract(b []byte, root, _ string) (e []library.Chunk) {
	r := ar.NewReader(bytes.NewReader(b))

	for {
		h, err := r.Next()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			slog.Error(err.Error())
			break
		}

		if h.Mode&int64(fs.ModeDir) != 0 {
			continue
		}

		buf, err := library.ReadMax(r, len(b))

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		e = append(e, library.Chunk{
			Path: sys.JoinPart(root, h.Name),
			Data: buf,
		})
	}

	return
}
