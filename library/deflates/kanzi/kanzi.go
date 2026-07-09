package kanzi

import (
	"bytes"
	"io"
	"log/slog"

	kio "github.com/flanglet/kanzi-go/v2/io"
	"go.foxforensics.eu/fox/v4/library"
)

func Detect(b []byte) bool {
	return library.HasMagic(b, 0, []byte{
		0x4B, 0x41, 0x4E, 0x5A,
	})
}

func Deflate(b []byte) ([]byte, error) {
	r, err := kio.NewReader(io.NopCloser(bytes.NewReader(b)), 4)

	if err != nil {
		return b, err
	}

	defer func() {
		if err := r.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()

	return library.ReadMax(r, len(b))
}
