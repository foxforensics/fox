package impfuzzy

import (
	"hash"
	"log"
	"strings"

	"github.com/glaslos/ssdeep"

	"go.foxforensics.dev/fox/v4/internal/pkg/hash/fuzzy"
)

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
	h.buf, err = fuzzy.GetImports(b, false)

	return len(b), err
}

func (h *ImpFuzzy) Sum(_ []byte) []byte {
	sum, err := ssdeep.FuzzyBytes([]byte(strings.Join(h.buf, ",")))

	if err != nil {
		log.Println(err)
	}

	return []byte(sum)
}
