package types

import (
	"bytes"
	"errors"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2/v2"

	"go.foxforensics.eu/fox/v4/internal/pkg/types/mmap"
)

const CR = '\n'

var re = regexp2.MustCompile(`^-?[0-9a-f]+[bhl]?$`)

type Limits struct {
	IsHead bool // is head limit
	IsTail bool // is tail limit
	Bytes  uint // bytes count
	Lines  uint // lines count
	Values struct {
		Bytes int
		Lines int
	}
}

func NewLimits(s string) (*Limits, error) {
	s = strings.TrimSpace(strings.ToLower(s))

	neg := strings.HasPrefix(s, "-")

	limits := &Limits{
		IsHead: !neg,
		IsTail: neg,
	}

	if len(s) == 0 {
		return limits, nil // empty
	}

	if ok, _ := re.MatchString(s); !ok {
		return nil, errors.New("invalid limit syntax")
	}

	var val int64
	var err error

	switch {
	case strings.HasSuffix(s, "h"):
		val, err = strconv.ParseInt(s[:len(s)-1], 16, 0)
	case strings.HasSuffix(s, "b"):
		val, err = strconv.ParseInt(s[:len(s)-1], 10, 0)
	case strings.HasSuffix(s, "l"):
		val, err = strconv.ParseInt(s[:len(s)-1], 10, 0)
	default:
		val, err = strconv.ParseInt(s, 10, 0)
	}

	if err != nil {
		return nil, err
	}

	if val < 0 {
		val *= -1
	}

	switch {
	case s[len(s)-1] == 'l':
		limits.Lines = uint(val)
	default:
		limits.Bytes = uint(val)
	}

	return limits, nil
}

func (l *Limits) Reduce(m mmap.MMap) mmap.MMap {
	var a, b = 0, len(m)

	if !l.IsHead && !l.IsTail {
		return m
	}

	if l.IsHead && l.Bytes > 0 {
		b = min(int(l.Bytes), b)
	}

	if l.IsTail && l.Bytes > 0 {
		a = max(len(m)-int(l.Bytes), 0)

		l.Values.Bytes = a
		l.Values.Lines = count(m) - count(m[a:])
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

		l.Values.Bytes = a
		l.Values.Lines = count(m) - n
	}

	return m[a:b]
}

func count(m mmap.MMap) int {
	v := bytes.Count(m, []byte{CR})

	if m[len(m)-1] != CR {
		v++ // last line
	}

	return v
}
