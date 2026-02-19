package ewf

import (
	"io"

	"github.com/Velocidex/go-ewf/parser"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte{
		'E', 'V', 'F', 0x09, 0x0D, 0x0A, 0xFF, 0x00, // EWF-E01, EWF-S01
	})
}

func Reader(f types.File) (io.ReaderAt, error) {
	r, err := parser.OpenEWFFile(nil, f)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func Combine(f ...types.File) (io.ReaderAt, error) {
	var rs []io.ReaderAt

	for _, r := range f {
		rs = append(rs, r)
	}

	r, err := parser.OpenEWFFile(nil, rs...)

	if err != nil {
		return nil, err
	}

	return r, nil
}
