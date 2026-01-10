package text

import (
	"fmt"
	"slices"
	"strings"
)

type Options struct {
	Min     uint
	Max     uint
	Sort    bool
	Wtf     int
	Find    []string
	First   bool
	Profile int
}

type String struct {
	org int
	Off string
	Cls string
	Str string
}

type Carver struct {
	opts     *Options
	cache    []*String
	strings  chan *String
	patterns *patterns
}

func NewCarver(opts *Options) *Carver {
	return &Carver{
		opts:     opts,
		cache:    make([]*String, 0),
		strings:  make(chan *String, opts.Profile*64),
		patterns: newPatterns(opts.Wtf),
	}
}

func (crv *Carver) Carve(block []byte) <-chan *String {
	go func() {
		var b []byte
		var c byte
		var o int

		for o, c = range block {
			if c >= SP && c <= DEL {
				b = append(b, c)
			} else {
				crv.flush(o, b)
				b = b[:0]
			}
		}

		crv.flush(o, b)
		close(crv.strings)
	}()

	if crv.opts.Sort {
		return crv.sort()
	} else {
		return crv.strings
	}
}

func (crv *Carver) flush(offset int, buf []byte) {
	str := string(buf)

	v := uint(len(strings.TrimSpace(str)))

	if v >= crv.opts.Min && v <= crv.opts.Max {
		org := max(offset-(len(buf)+1), 0)
		off := fmt.Sprintf("%08x", org)
		cls := ""

		if crv.opts.Wtf > 0 {
			res := crv.patterns.Search(str)

			if len(crv.opts.Find) > 0 && !res.Match(crv.opts.Find) {
				return
			}

			cls = res.ToString(crv.opts.First)
		}

		crv.strings <- &String{org, off, cls, str}
	}
}

func (crv *Carver) sort() <-chan *String {
	sorted := make(chan *String, cap(crv.strings))

	go func() {
		for s := range crv.strings {
			crv.cache = append(crv.cache, s)
		}

		slices.SortStableFunc(crv.cache, func(a, b *String) int {
			if a.Str != b.Str {
				return strings.Compare(a.Str, b.Str)
			} else {
				return a.org - b.org
			}
		})

		for _, s := range crv.cache {
			sorted <- s
		}

		close(sorted)
	}()

	return sorted
}
