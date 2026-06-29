package cpio

import (
	"bytes"
	"errors"
	"io"
	"log/slog"

	"github.com/cavaliergopher/cpio"
	"go.foxforensics.eu/fox/v4/internal/pkg/lib"
	"go.foxforensics.eu/fox/v4/internal/sys"
)

func Detect(b []byte) bool {
	for _, m := range [][]byte{
		{0x30, 0x37, 0x30, 0x37, 0x30, 0x31}, // SRV4
		{0x30, 0x37, 0x30, 0x37, 0x30, 0x32}, // SRV4 with CRC
	} {
		if lib.HasMagic(b, 0, m) {
			return true
		}
	}

	return false
}

func Extract(b []byte, root, _ string) (e []lib.Stream) {
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
