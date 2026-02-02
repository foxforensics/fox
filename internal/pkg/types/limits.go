package types

import (
	"bytes"
	"log"
	"strconv"
	"strings"

	"github.com/cuhsat/fox/v4/internal/pkg/types/mmap"
)

const CR = '\n'

type Limits struct {
	IsHead bool // is head limit
	IsTail bool // is tail limit
	Bytes  uint // bytes count
	Lines  uint // lines count
	Offset struct {
		Bytes int
		Lines int
	}
}

func NewLimits(h, t bool, b, l string) *Limits {
	return &Limits{
		IsHead: h,
		IsTail: t,
		Bytes:  convert(b),
		Lines:  convert(l),
	}
}

func (l *Limits) Reduce(m mmap.MMap) mmap.MMap {
	var a, b = 0, len(m)

	if l.IsHead && l.Bytes > 0 {
		b = min(int(l.Bytes), b)
	}

	if l.IsTail && l.Bytes > 0 {
		a = max(len(m)-int(l.Bytes), 0)

		l.Offset.Bytes = a
		l.Offset.Lines = count(m) - count(m[a:])
	}

	if l.IsHead && l.Lines > 0 {
		i := a

		for n := 0; i < b && n < int(l.Lines); i++ {
			if m[i] == CR {
				n++
			}
		}

		b = min(i, b)
	}

	if l.IsTail && l.Lines > 0 {
		i, n := b-1, 0

		for ; i > a && n < int(l.Lines); i-- {
			if m[i-1] == CR {
				n++
			}
		}

		a = max(i, a)

		if a > 0 {
			a++ // skip linebreak
		}

		if i == 0 {
			n++ // add first line
		}

		l.Offset.Bytes = a
		l.Offset.Lines = count(m) - n
	}

	return m[a:b]
}

func convert(s string) uint {
	var val uint64
	var err error

	if len(s) == 0 {
		return 0
	}

	s = strings.ToLower(s)

	switch {
	case strings.HasPrefix(s, "0x"):
		val, err = strconv.ParseUint(s[2:], 16, 0)
	case strings.HasPrefix(s, "#"):
		val, err = strconv.ParseUint(s[1:], 16, 0)
	default:
		val, err = strconv.ParseUint(s, 10, 0)
	}

	if err != nil {
		log.Fatalln(err)
	}

	return uint(val)
}

func count(m mmap.MMap) int {
	v := bytes.Count(m, []byte{CR})

	if m[len(m)-1] != CR {
		v++
	}

	return v
}
