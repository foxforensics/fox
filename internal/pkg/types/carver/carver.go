package carver

import (
	"fmt"
	"slices"
	"strings"

	"go.foxforensics.dev/ustrings/ustrings"
)

type Options struct {
	Min      uint
	Max      uint
	Ascii    bool
	Sort     bool
	Wtf      int
	Find     []string
	First    bool
	Parallel int
}

type String struct {
	ustrings.String
	Address string
	Classes string
}

type Carver struct {
	opts  *Options
	cache []*String
	ch    chan *String
	db    Database
}

func New(opts *Options) *Carver {
	return &Carver{
		opts:  opts,
		cache: make([]*String, 0),
		ch:    make(chan *String, opts.Parallel*64),
		db:    buildDB(opts.Wtf),
	}
}

func (crv *Carver) Carve(block []byte) <-chan *String {
	go func() {
		defer close(crv.ch)

		for str := range ustrings.Carve(
			block,
			crv.opts.Min,
			crv.opts.Max,
			true,
			crv.opts.Ascii,
		) {
			var adr = fmt.Sprintf("%08x", str.Offset)
			var cls string

			// append class
			if crv.opts.Wtf > 0 {
				v := crv.db.Lookup(str.Value)

				// search classes
				if len(crv.opts.Find) > 0 && !contains(v, crv.opts.Find) {
					continue
				}

				// format classes
				if !crv.opts.First {
					cls = strings.Join(v, " ")
				} else if len(v) > 0 {
					cls = v[0]
				}
			}

			crv.ch <- &String{*str, adr, cls}
		}
	}()

	if crv.opts.Sort {
		return crv.sort()
	}

	return crv.ch
}

func (crv *Carver) sort() <-chan *String {
	sorted := make(chan *String, cap(crv.ch))

	go func() {
		defer close(sorted)

		for s := range crv.ch {
			crv.cache = append(crv.cache, s)
		}

		slices.SortStableFunc(crv.cache, func(a, b *String) int {
			if a.Value != b.Value {
				return strings.Compare(a.Value, b.Value)
			}

			return int(a.Offset - b.Offset)
		})

		for _, s := range crv.cache {
			sorted <- s
		}
	}()

	return sorted
}

func contains(a, b []string) bool {
	for _, x := range a {
		for _, y := range b {
			if strings.Compare(strings.ToLower(x), strings.ToLower(y)) == 0 {
				return true
			}
		}
	}

	return false
}
