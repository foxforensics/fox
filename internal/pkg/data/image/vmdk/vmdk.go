package vmdk

import (
	"bytes"
	"io"

	"github.com/Velocidex/go-vmdk/parser"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte{
		'K', 'D', 'M', 'V', // VMDK Multi-Extent sparse
	})
}

func Ingest(b []byte) ([]byte, error) {
	vol, err := parser.GetVMDKContext(bytes.NewReader(b), len(b), func(filename string) (reader io.ReaderAt, closer func(), err error) {
		return bytes.NewReader(b), func() {}, nil
	})

	if err != nil {
		return b, err
	}

	buf := make([]byte, vol.Size())

	if _, err = vol.ReadAt(buf, 0); err != nil {
		return b, err
	}

	return buf, nil
}
