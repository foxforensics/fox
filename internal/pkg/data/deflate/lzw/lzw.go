package lzw

import (
	"bytes"
	"compress/lzw"
	"errors"
	"io"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte{
		0x1F, 0x9D,
	})
}

func Deflate(b []byte) ([]byte, error) {
	// compress compatible settings
	r := lzw.NewReader(bytes.NewReader(b[3:]), lzw.LSB, 8)

	defer func(r io.Closer) {
		_ = r.Close()
	}(r)

	buf, err := io.ReadAll(r)

	// ignore errors for faulty compress files
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, err
	}

	return buf, nil
}
