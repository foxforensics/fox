package carver

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"

	db "go.foxforensics.eu/fox/v4/internal/pkg/str"
	fstrings "go.foxforensics.eu/strings/strings"
)

var (
	entries db.Database
	buildDB sync.Once
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
}

type Carver struct {
	opts    *Options
	entries db.Database
}

func New(opts *Options) *Carver {
	buildDB.Do(func() {
		entries = db.BuildDB(opts.What)
	})

	return &Carver{
		opts:    opts,
		entries: entries,
	}
}

func (crv *Carver) Carve(ctx context.Context, block []byte) <-chan *String {
	ch := make(chan *String, crv.opts.Threads*64)

	go func() {
		defer close(ch)

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
			case ch <- &String{*str, adr, cls}:
			case <-ctx.Done():
				return
			}
		}
	}()

	if crv.opts.Sort {
		return crv.sort(ctx, ch)
	}

	return ch
}

func (crv *Carver) sort(ctx context.Context, ch <-chan *String) <-chan *String {
	sorted := make(chan *String, cap(ch))

	go func() {
		defer close(sorted)

		v := make([]String, 0)

		for s := range ch {
			v = append(v, *s)
		}

		slices.SortStableFunc(v, compare)

		for _, s := range v {
			select {
			case sorted <- &s:
			case <-ctx.Done():
				return
			}
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
