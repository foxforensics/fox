package bgzf

import (
	"bytes"
	"io"

	"go.foxforensics.dev/go-bgzf/bgzf"

	"go.foxforensics.dev/fox/v4/internal/pkg/file"
)

func Detect(b []byte) bool {
	return file.HasMagic(b, 0, []byte{
		0x1F, 0x8B, 0x08, 0x04,
	})
}

func Deflate(b []byte) ([]byte, error) {
	r, err := bgzf.NewReader(bytes.NewReader(b), 0)

	if err != nil {
		return b, err
	}

	defer func(r *bgzf.Reader) {
		_ = r.Close()
	}(r)

	return io.ReadAll(r)
}
