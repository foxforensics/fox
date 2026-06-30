package ar

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"log/slog"

	"github.com/mkrautz/goar"
	"go.foxforensics.eu/fox/v4/internal/lib"
	"go.foxforensics.eu/fox/v4/internal/sys"
)

func Detect(b []byte) bool {
	return lib.HasMagic(b, 0, []byte{
		0x21, 0x3C, 0x61, 0x72, 0x63, 0x68, 0x3E, 0x0A,
	})
}

func Extract(b []byte, root, _ string) (e []lib.Stream) {
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

		buf, err := io.ReadAll(r)

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		e = append(e, lib.Stream{
			Path: sys.JoinPart(root, h.Name),
			Data: buf,
		})
	}

	return
}
