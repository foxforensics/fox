package vhdx

import (
	"bytes"

	"github.com/Velocidex/go-vhdx/parser"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte{
		'v', 'h', 'd', 'x', 'f', 'i', 'l', 'e', // Virtual Hard Disk v2
	})
}

func Ingest(b []byte) ([]byte, error) {
	vol, err := parser.NewVHDXFile(bytes.NewReader(b))

	if err != nil {
		return b, err
	}

	buf := make([]byte, vol.Metadata.VirtualDiskSize)

	if _, err = vol.ReadAt(buf, 0); err != nil {
		return b, err
	}

	return buf, nil
}
