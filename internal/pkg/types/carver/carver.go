package carver

import (
	"fmt"
	"log"
	"slices"
	"strings"

	fstrings "go.foxforensics.dev/strings/strings"

	"go.foxforensics.dev/checker/services"
	"go.foxforensics.dev/checker/services/vt"

	"go.foxforensics.dev/fox/v4/internal/pkg/text"
)

type Options struct {
	Min      uint
	Max      uint
	Ascii    bool
	Sort     bool
	Trim     bool
	What     int
	Find     []string
	First    bool
	Lookup   bool
	Parallel int
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
	entries text.Database
}

func New(opts *Options) *Carver {
	return &Carver{
		opts:    opts,
		list:    make([]String, 0),
		strings: make(chan *String, opts.Parallel*64),
		entries: text.BuildDB(opts.What),
	}
}

func (crv *Carver) Carve(block []byte) <-chan *String {
	go func() {
		defer close(crv.strings)

		for str := range fstrings.Carve(
			block,
			crv.opts.Min,
			crv.opts.Max,
			crv.opts.Ascii,
			crv.opts.Trim,
		) {
			var adr = fmt.Sprintf("%08x", str.Offset)
			var cls string
			var sus bool

			// lookup classes
			if crv.opts.What > 0 {
				v := crv.entries.Lookup(str.Value)

				// search entries
				if len(crv.opts.Find) > 0 && !contains(v, crv.opts.Find) {
					continue
				}

				// format entries
				if !crv.opts.First {
					cls = strings.Join(v, " ")
				} else if len(v) > 0 {
					cls = v[0]
				}

				// check entries
				if crv.opts.Lookup {
					sus = check(str.Value, v)
				}
			}

			crv.strings <- &String{*str, adr, cls, sus}
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

func check(s string, v []string) bool {
	var res *services.Result
	var err error

	switch {
	case slices.Contains(v, "IPv6"):
		fallthrough
	case slices.Contains(v, "IPv4"):
		res, err = vt.CheckIp(s)
	case slices.Contains(v, "URL"):
		res, err = vt.CheckUrl(s)
	}

	if err != nil {
		log.Printf("warning: %s!\n", err.Error())
		return false
	}

	if res == nil {
		return false
	}

	return res.Verdict == services.Suspicious
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
		return int(a.Offset - b.Offset)
	}

	return strings.Compare(a.Value, b.Value)
}
