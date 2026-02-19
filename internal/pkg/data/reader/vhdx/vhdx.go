package vhdx

import (
	"io"

	"github.com/Velocidex/go-vhdx/parser"

	"foxhunt.dev/fox/internal/pkg/data"
	"foxhunt.dev/fox/internal/pkg/types"
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte{
		'v', 'h', 'd', 'x', 'f', 'i', 'l', 'e', // Virtual Hard Disk v2
	})
}

func Reader(f types.File) (io.ReaderAt, error) {
	r, err := parser.NewVHDXFile(f)

	if err != nil {
		return nil, err
	}

	return r, nil
}
