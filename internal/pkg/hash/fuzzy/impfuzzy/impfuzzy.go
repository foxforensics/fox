package impfuzzy

import (
	"bytes"
	"debug/pe"
	"errors"
	"fmt"
	"hash"
	"log"
	"strings"

	"github.com/glaslos/ssdeep"

	intern "foxhunt.dev/fox/internal/pkg/data/convert/bin/pe"
)

var ErrNotSupported = errors.New("file type not supported")

type ImpFuzzy struct {
	buf []string
}

func New() hash.Hash {
	return new(ImpFuzzy)
}

func (h *ImpFuzzy) BlockSize() int {
	return h.BlockSize()
}

func (h *ImpFuzzy) Size() int {
	return h.Size()
}

func (h *ImpFuzzy) Reset() {
	h.buf = h.buf[:0]
}

func (h *ImpFuzzy) Write(b []byte) (n int, err error) {
	if !intern.Detect(b) {
		return 0, ErrNotSupported
	}

	f, err := pe.NewFile(bytes.NewReader(b))

	if err != nil {
		return 0, err
	}

	defer func(f *pe.File) {
		_ = f.Close()
	}(f)

	iat, err := f.ImportedSymbols()

	if err != nil {
		return 0, err
	}

	rep := strings.NewReplacer(".dll", "", ".ocx", "", ".sys", "")

	for _, e := range iat {
		if !strings.Contains(e, ":") {
			continue
		}

		p := strings.Split(rep.Replace(strings.ToLower(e)), ":")

		h.buf = append(h.buf, fmt.Sprintf("%s.%s", p[1], p[0]))
	}

	return len(b), nil
}

func (h *ImpFuzzy) Sum(_ []byte) []byte {
	sum, err := ssdeep.FuzzyBytes([]byte(strings.Join(h.buf, ",")))

	if err != nil {
		log.Println(err)
	}

	return []byte(sum)
}
