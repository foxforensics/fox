package types

import "github.com/cuhsat/go-mmap"

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

func (l *Limits) Reduce(m mmap.MMap) mmap.MMap {
	var a, b = 0, len(m)

	if l.IsHead && l.Bytes > 0 {
		b = min(int(l.Bytes), b)
	}

	if l.IsTail && l.Bytes > 0 {
		a = max(len(m)-int(l.Bytes), 0)
		l.Offset.Bytes = a
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

		l.Offset.Lines = n
	}

	return m[a:b]
}
