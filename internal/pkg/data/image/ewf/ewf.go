package ewf

import (
	"bytes"

	"github.com/Velocidex/go-ewf/parser"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte{
		'E', 'V', 'F', 0x09, 0x0D, 0x0A, 0xFF, 0x00, // EWF-E01, EWF-S01
	})
}

func Ingest(b []byte) ([]byte, error) {
	vol, err := parser.OpenEWFFile(nil, bytes.NewReader(b))

	if err != nil {
		return b, err
	}

	buf := make([]byte, vol.ChunkSize*vol.NumberOfChunks)

	if _, err = vol.ReadAt(buf, 0); err != nil {
		return b, err
	}

	return buf, nil
}
