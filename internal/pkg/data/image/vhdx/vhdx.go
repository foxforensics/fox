package vhdx

import (
	"bytes"

	"github.com/Velocidex/go-vhdx/parser"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
	"github.com/cuhsat/fox/v4/internal/pkg/types/mmap"
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

	buf := mmap.Remap(vol, int(vol.Metadata.VirtualDiskSize))

	mmap.Unmap(b)

	return buf, nil
}
