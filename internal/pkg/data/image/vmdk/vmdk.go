package vmdk

import (
	"bytes"
	"io"

	"github.com/Velocidex/go-vmdk/parser"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
	"github.com/cuhsat/fox/v4/internal/pkg/types/mmap"
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

	buf := mmap.Remap(vol, int(vol.Size()))

	mmap.Unmap(b)

	return buf, nil
}
