package carver

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strings"

	fstrings "go.foxforensics.eu/strings/strings"
)

type Options struct {
	Min     uint
	Max     uint
	Ascii   bool
	Sort    bool
	Trim    bool
	What    int
	Class   []string
	Threads int
}

type String struct {
	fstrings.String
	Address string
	Classes string
	Suspect bool
}

type Carver struct {
	opts    *Options
	list    []String
	strings chan *String
	entries Database
}

func New(opts *Options) *Carver {
	return &Carver{
		opts:    opts,
		list:    make([]String, 0),
		strings: make(chan *String, opts.Threads*64),
		entries: BuildDB(opts.What),
	}
}

func (crv *Carver) Carve(ctx context.Context, block []byte) <-chan *String {
	go func() {
		defer close(crv.strings)

		for str := range fstrings.Carve(
			ctx,
			block,
			crv.opts.Min,
			crv.opts.Max,
			crv.opts.Ascii,
			crv.opts.Trim,
		) {
			var adr = fmt.Sprintf("%08x", str.Offset)
			var cls string

			// lookup classes
			if crv.opts.What > 0 {
				v := crv.entries.Lookup(str.Value)

				// search entries
				if len(crv.opts.Class) > 0 && !contains(v, crv.opts.Class) {
					continue
				}

				// format entries
				cls = strings.Join(v, " ")
			}

			select {
			case <-ctx.Done():
				return
			default:
				crv.strings <- &String{*str, adr, cls, false}
			}
		}
	}()

	if crv.opts.Sort {
		return crv.sort()
	}

	return crv.strings
}

func (crv *Carver) sort() <-chan *String {
	sorted := make(chan *String, cap(crv.strings))

	go func() {
		defer close(sorted)

		for s := range crv.strings {
			crv.list = append(crv.list, *s)
		}

		slices.SortStableFunc(crv.list, compare)

		for _, s := range crv.list {
			sorted <- &s
		}
	}()

	return sorted
}

func contains(a, b []string) bool {
	for _, x := range a {
		for _, y := range b {
			if strings.EqualFold(x, y) {
				return true
			}
		}
	}

	return false
}

func compare(a, b String) int {
	if a.Value == b.Value {
		return cmp.Compare(a.Offset, b.Offset)
	}

	return strings.Compare(a.Value, b.Value)
}
