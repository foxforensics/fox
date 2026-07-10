package rar

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"strings"

	"github.com/nwaples/rardecode/v2"
	"go.foxforensics.eu/fox/v5/internal/sys"
	"go.foxforensics.eu/fox/v5/library"
	"go.foxforensics.eu/fox/v5/library/archives"
)

func Detect(b []byte) bool {
	return library.HasMagic(b, 0, []byte{
		0x52, 0x61, 0x72, 0x21, 0x1A, 0x07,
	})
}

func Extract(b []byte, root, pass string) (e []library.Chunk) {
	r, err := rardecode.NewReader(bytes.NewBuffer(b), rardecode.Password(pass))

	if err != nil {
		slog.Error(err.Error())
		return
	}

	for {
		h, err := r.Next()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			slog.Error(err.Error())
			break
		}

		if strings.HasSuffix(h.Name, "/") {
			continue
		}

		buf, err := library.ReadMax(r, len(b))

		if err != nil {
			slog.Error(archives.ErrCorruptPassword.Error())
			continue
		}

		e = append(e, library.Chunk{
			Path: sys.JoinPart(root, h.Name),
			Data: buf,
		})
	}

	return
}
