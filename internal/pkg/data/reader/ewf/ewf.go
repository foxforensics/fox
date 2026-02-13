package ewf

import (
	"io"

	"github.com/Velocidex/go-ewf/parser"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte{
		'E', 'V', 'F', 0x09, 0x0D, 0x0A, 0xFF, 0x00, // EWF-E01, EWF-S01
	})
}

func Reader(r ...io.ReaderAt) (io.ReaderAt, error) {
	v, err := parser.OpenEWFFile(nil, r...)

	if err != nil {
		return nil, err
	}

	return v, nil
}
