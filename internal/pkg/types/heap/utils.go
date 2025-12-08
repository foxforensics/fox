package heap

import (
	"fmt"
	"math"
	"strings"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
)

type String struct {
	Off string
	Cls string
	Str string
}

func (h *Heap) Entropy(mn, mx float64) (float64, bool) {
	var a [256]float64
	var v float64

	for _, b := range h.mmap {
		a[b]++
	}

	l := float64(len(h.MMap()))

	for i := range 256 {
		if a[i] != 0 {
			f := a[i] / l
			v -= f * math.Log2(f)
		}
	}

	v /= 8

	// heap filtered
	if v < mn || v > mx {
		return 0, false
	}

	return v, true
}

func (h *Heap) Strings(mn, mx uint, wtf int, fst bool) <-chan String {
	var ch = make(chan String, 4096)
	var db *text.Strings
	var buf []byte
	var off int
	var b byte

	if wtf > 0 {
		db = text.GetStrings(wtf)
	}

	// flush closure
	flush := func() {
		cls, str := "", string(buf)

		v := uint(len(strings.TrimSpace(str)))

		if v >= mn && v <= mx {
			if db != nil {
				cls = db.Search(str).ToString(fst)
			}

			ch <- String{
				fmt.Sprintf("%08x", max(off-(len(buf)+1), 0)),
				cls,
				str,
			}
		}

		buf = buf[:0]
	}

	// carve closure
	carve := func() {
		for off, b = range h.mmap {
			if b >= text.SP && b <= text.DEL {
				buf = append(buf, b)
			} else {
				flush()
			}
		}

		flush()
		close(ch)
	}

	go carve()

	return ch
}
