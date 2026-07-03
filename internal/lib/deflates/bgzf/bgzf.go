package bgzf

import (
	"bytes"
	"log/slog"

	"go.foxforensics.eu/fox/v4/internal/lib"
	"go.foxforensics.eu/go-bgzf/bgzf"
)

func Detect(b []byte) bool {
	return lib.HasMagic(b, 0, []byte{
		0x1F, 0x8B, 0x08, 0x04,
	})
}

func Deflate(b []byte) ([]byte, error) {
	r, err := bgzf.NewReader(bytes.NewReader(b), 0)

	if err != nil {
		return b, err
	}

	defer func() {
		if err := r.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()

	return lib.ReadMax(r, len(b))
}
