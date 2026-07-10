package cpio

import (
	"bytes"
	"errors"
	"io"
	"log/slog"

	"github.com/cavaliergopher/cpio"
	"go.foxforensics.eu/fox/v5/internal/sys"
	"go.foxforensics.eu/fox/v5/library"
)

func Detect(b []byte) bool {
	for _, m := range [][]byte{
		{0x30, 0x37, 0x30, 0x37, 0x30, 0x31}, // SRV4
		{0x30, 0x37, 0x30, 0x37, 0x30, 0x32}, // SRV4 with CRC
	} {
		if library.HasMagic(b, 0, m) {
			return true
		}
	}

	return false
}

func Extract(b []byte, root, _ string) (e []library.Chunk) {
	r := cpio.NewReader(bytes.NewBuffer(b))

	for {
		h, err := r.Next()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			slog.Error(err.Error())
			break
		}

		if !h.Mode.IsRegular() {
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
