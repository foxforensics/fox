package vmdk

import (
	"io"
	"os"

	"github.com/Velocidex/go-vmdk/parser"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

const size = 64 * 1024 // buffer size

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte{
		'K', 'D', 'M', 'V', // VMDK Multi-Extent sparse
	})
}

func Reader(f *os.File) (io.ReaderAt, error) {
	r, err := parser.GetVMDKContext(f, size, func(filename string) (io.ReaderAt, func(), error) {
		return f, func() { _ = f.Close() }, nil
	})

	if err != nil {
		return nil, err
	}

	return r, nil
}
