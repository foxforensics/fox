package kanzi

import (
	"bytes"
	"io"

	kio "github.com/flanglet/kanzi-go/v2/io"

	"go.foxforensics.dev/fox/v4/internal/pkg/file"
)

func Detect(b []byte) bool {
	return file.HasMagic(b, 0, []byte{
		0x4B, 0x41, 0x4E, 0x5A,
	})
}

func Deflate(b []byte) ([]byte, error) {
	r, err := kio.NewReader(io.NopCloser(bytes.NewReader(b)), 4)

	if err != nil {
		return b, err
	}

	defer func(r io.Closer) {
		_ = r.Close()
	}(r)

	return io.ReadAll(r)
}
