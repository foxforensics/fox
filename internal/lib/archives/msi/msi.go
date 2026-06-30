package msi

import (
	"bytes"
	"io"
	"log/slog"

	"go.foxforensics.eu/fox/v4/internal/lib"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/go-msi/msi"
)

func Detect(b []byte) bool {
	return lib.HasMagic(b, 0, []byte{
		0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1,
	})
}

func Extract(b []byte, root, _ string) (e []lib.Stream) {
	r, err := msi.Open(bytes.NewReader(b))

	if err != nil {
		slog.Error(err.Error())
		return
	}

	streams := r.Streams()

	for {
		name := streams.Next()

		if len(name) == 0 {
			break
		}

		str, err := r.ReadStream(name)

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		buf, err := io.ReadAll(str)

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		e = append(e, lib.Stream{
			Path: sys.JoinPart(root, name),
			Data: buf,
		})
	}

	return
}
