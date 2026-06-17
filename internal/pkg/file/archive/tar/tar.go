package tar

import (
	"archive/tar"
	"bytes"
	"io"
	"log/slog"
	"strings"

	"go.foxforensics.eu/fox/v4/internal/pkg/file"
)

func Detect(b []byte) bool {
	return file.HasMagic(b, 257, []byte{
		0x75, 0x73, 0x74, 0x61, 0x72,
	})
}

func Extract(b []byte, root, _ string) (e []file.Stream) {
	r := tar.NewReader(bytes.NewBuffer(b))

	for {
		h, err := r.Next()

		if err == io.EOF {
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

		e = append(e, file.Stream{
			Path: file.JoinPart(root, h.Name),
			Data: buf,
		})
	}

	return
}
