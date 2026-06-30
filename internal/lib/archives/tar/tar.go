package tar

import (
	"archive/tar"
	"bytes"
	"errors"
	"io"
	"log/slog"
	"strings"

	"go.foxforensics.eu/fox/v4/internal/lib"
	"go.foxforensics.eu/fox/v4/internal/sys"
)

func Detect(b []byte) bool {
	return lib.HasMagic(b, 257, []byte{
		0x75, 0x73, 0x74, 0x61, 0x72,
	})
}

func Extract(b []byte, root, _ string) (e []lib.Stream) {
	r := tar.NewReader(bytes.NewBuffer(b))

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
