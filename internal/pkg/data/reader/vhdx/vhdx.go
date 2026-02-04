package vhdx

import (
	"io"
	"os"

	"github.com/Velocidex/go-vhdx/parser"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte{
		'v', 'h', 'd', 'x', 'f', 'i', 'l', 'e', // Virtual Hard Disk v2
	})
}

func Reader(f *os.File) (io.ReaderAt, error) {
	r, err := parser.NewVHDXFile(f)

	if err != nil {
		return nil, err
	}

	return r, nil
}
