package carver

import (
	"fmt"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Options struct {
	Min     uint
	Max     uint
	Ascii   bool
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
	stream := make(chan byte)

	go func() {
		for _, b := range block {
			stream <- b
		}

		close(stream)
	}()

	go func() {
		var off, i, l int
		var buf []rune

		cp := make([]byte, 4)

		for b := range stream {
			cp[0], off, i = b, off+1, 1

			// fill remaining bytes
			if !cvr.opts.Ascii {
				for ; i < bytes(b); i++ {
					if b, ok := <-stream; ok {
						cp[i], off = b, off+1
					} else {
						break
					}
				}
			}

			// convert code point to rune
			r, _ := utf8.DecodeRune(cp[:i])

			// append rune to buffer or flush
			if r != utf8.RuneError && unicode.IsPrint(r) {
				buf, l = append(buf, r), l+i
			} else if len(buf) > 0 {
				cvr.flush(max(off-l-1, 0), buf)
				buf, l = buf[:0], 0
			}
		}

		cvr.flush(max(off-l-1, 0), buf)

		close(cvr.ch)
	}()

	if cvr.opts.Sort {
		return cvr.sort()
	}

	return cvr.ch
}

func (cvr *Carver) flush(off int, buf []rune) {
	str := string(buf)

	v := uint(len(strings.TrimSpace(str)))

	if v >= cvr.opts.Min && v <= cvr.opts.Max {
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
			} else if len(v) > 0 {
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
			}

			return a.off - b.off
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

func bytes(b byte) int {
	switch {
	default:
		return 1
	case b&0x80 == 0:
		return 1
	case b&0xE0 == 0xC0:
		return 2
	case b&0xF0 == 0xE0:
		return 3
	case b&0xF8 == 0xF0:
		return 4
	}
}
