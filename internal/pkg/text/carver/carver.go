package carver

import (
	"fmt"
	"slices"
	"strings"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
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
	off int
	Adr string
	Cls string
	Str string
}

type Carver struct {
	opts  *Options
	cache []*String
	ch    chan *String
	db    database
}

func New(opts *Options) *Carver {
	return &Carver{
		opts:  opts,
		cache: make([]*String, 0),
		ch:    make(chan *String, opts.Profile*64),
		db:    buildDB(opts.Wtf),
	}
}

func (cvr *Carver) Carve(block []byte) <-chan *String {
	go func() {
		var b []byte
		var c byte
		var o int

		for o, c = range block {
			if c >= text.SP && c <= text.DEL {
				b = append(b, c)
			} else {
				cvr.flush(o, b)
				b = b[:0]
			}
		}

		cvr.flush(o, b)
		close(cvr.ch)
	}()

	if cvr.opts.Sort {
		return cvr.sort()
	} else {
		return cvr.ch
	}
}

func (cvr *Carver) flush(off int, buf []byte) {
	str := string(buf)

	v := uint(len(strings.TrimSpace(str)))

	if v >= cvr.opts.Min && v <= cvr.opts.Max {
		off = max(off-(len(buf)+1), 0)
		adr := fmt.Sprintf("%08x", off)
		cls := ""

		if cvr.opts.Wtf > 0 {
			v := cvr.db.Search(str)

			// search classes
			if len(cvr.opts.Find) > 0 && !contains(v, cvr.opts.Find) {
				return
			}

			// format classes
			if !cvr.opts.First {
				cls = strings.Join(v, ", ")
			} else {
				cls = v[0]
			}
		}

		cvr.ch <- &String{off, adr, cls, str}
	}
}

func (cvr *Carver) sort() <-chan *String {
	sorted := make(chan *String, cap(cvr.ch))

	go func() {
		for s := range cvr.ch {
			cvr.cache = append(cvr.cache, s)
		}

		slices.SortStableFunc(cvr.cache, func(a, b *String) int {
			if a.Str != b.Str {
				return strings.Compare(a.Str, b.Str)
			} else {
				return a.off - b.off
			}
		})

		for _, s := range cvr.cache {
			sorted <- s
		}

		close(sorted)
	}()

	return sorted
}

func contains(a, b []string) bool {
	for _, x := range a {
		for _, y := range b {
			if strings.Compare(
				strings.ToLower(x),
				strings.ToLower(y),
			) == 0 {
				return true
			}
		}
	}

	return false
}
