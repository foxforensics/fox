package rar

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"strings"

	"github.com/nwaples/rardecode/v2"
	"go.foxforensics.eu/fox/v4/internal/pkg/lib"
	"go.foxforensics.eu/fox/v4/internal/pkg/lib/archives"
	"go.foxforensics.eu/fox/v4/internal/sys"
)

func Detect(b []byte) bool {
	return lib.HasMagic(b, 0, []byte{
		0x52, 0x61, 0x72, 0x21, 0x1A, 0x07,
	})
}

func Extract(b []byte, root, pass string) (e []lib.Stream) {
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

		buf, err := io.ReadAll(r)

		if err != nil {
			slog.Error(archives.ErrCorruptPassword.Error())
			continue
		}

		e = append(e, lib.Stream{
			Path: sys.JoinPart(root, h.Name),
			Data: buf,
		})
	}

	return
}
