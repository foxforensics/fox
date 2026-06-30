package lzw

import (
	"bytes"
	"compress/lzw"
	"errors"
	"io"
	"log/slog"

	"go.foxforensics.eu/fox/v4/internal/lib"
)

func Detect(b []byte) bool {
	return lib.HasMagic(b, 0, []byte{
		0x1F, 0x9D,
	})
}

func Deflate(b []byte) ([]byte, error) {
	if len(b) < 4 {
		return b, errors.New("invalid length")
	}

	// compress compatible settings
	r := lzw.NewReader(bytes.NewReader(b[3:]), lzw.LSB, 8)

	defer func() {
		if err := r.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()

	buf, err := io.ReadAll(r)

	// ignore errors for faulty compress files
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		return buf, err
	}

	return buf, nil
}
