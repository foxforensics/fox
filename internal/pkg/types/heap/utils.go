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

func (h *Heap) Strings(m, x uint, w int, s []string, f bool, p int) <-chan String {
	var ch = make(chan String, p*64)
	var db *text.Strings
	var buf []byte
	var off int
	var b byte

	if w > 0 {
		db = text.GetStrings(w)
	}

	// flush closure
	flush := func() {
		defer func() {
			buf = buf[:0]
		}()

		str := string(buf)

		v := uint(len(strings.TrimSpace(str)))

		if v >= m && v <= x {
			off := fmt.Sprintf("%08x", max(off-(len(buf)+1), 0))
			cls := ""

			if db != nil {
				res := db.Search(str)

				if len(s) > 0 && !res.Match(s) {
					return
				}

				cls = res.ToString(f)
			}

			ch <- String{off, cls, str}
		}
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

func (h *Heap) Entropy(block []byte) float64 {
	var a [256]float64
	var v float64

	for _, b := range block {
		a[b]++
	}

	l := float64(len(block))

	for i := range 256 {
		if a[i] != 0 {
			f := a[i] / l
			v -= f * math.Log2(f)
		}
	}

	v /= 8

	return v
}
