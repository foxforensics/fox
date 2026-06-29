package pkg

import (
	"bytes"
	"errors"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2/v2"
	"go.foxforensics.eu/fox/v4/internal/pkg/cat/smap"
	"go.foxforensics.eu/fox/v4/internal/sys/mmap"
)

const CR = '\n'

var limit = regexp2.MustCompile(`^-?[0-9a-f]+[bhl]?$`)

type Query struct {
	Regex  *regexp2.Regexp // regex
	Bytes  uint            // bytes count
	Lines  uint            // lines count
	Before uint            // lines before
	After  uint            // lines after
	IsHead bool            // has head limit
	IsTail bool            // has tail limit
	Values struct {
		Bytes int
		Lines int
	}
}

func SetQuery(q *Query, s string, re *regexp2.Regexp) error {
	s = strings.TrimSpace(strings.ToLower(s))

	neg := strings.HasPrefix(s, "-")

	q.Regex = re
	q.IsHead = !neg
	q.IsTail = neg

	if len(s) == 0 {
		return nil // empty
	}

	if ok, _ := limit.MatchString(s); !ok {
		return errors.New("invalid limit syntax")
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
		return err
	}

	if val < 0 {
		val *= -1
	}

	switch {
	case s[len(s)-1] == 'l':
		q.Lines = uint(val)
	default:
		q.Bytes = uint(val)
	}

	return nil
}

func (q *Query) Filter(s smap.SMap, n int) smap.SMap {
	if q.Regex == nil {
		return s // not filtered
	}

	v := s.Grep(q.Regex, n)

	if q.Before+q.After == 0 {
		return v // without context
	}

	r := make(smap.SMap, 0, len(v))

	for grp, str := range v {
		for _, b := range (s)[max(int((str.Line-1)-q.Before), 0) : str.Line-1] {
			b.Group = uint(grp + 1) // modify copy
			r = append(r, b)
		}

		str.Group = uint(grp + 1)
		r = append(r, str)

		for _, a := range (s)[str.Line:min(int(str.Line+q.After), len(s))] {
			a.Group = uint(grp + 1) // modify copy
			r = append(r, a)
		}
	}

	return r // with context
}

func (q *Query) Reduce(m mmap.MMap) mmap.MMap {
	var a, b = 0, len(m)

	if !q.IsHead && !q.IsTail {
		return m
	}

	if q.IsHead && q.Bytes > 0 {
		b = min(int(q.Bytes), b)
	}

	if q.IsTail && q.Bytes > 0 {
		a = max(len(m)-int(q.Bytes), 0)

		// save last lines
		q.Values.Bytes = a
		q.Values.Lines = count(m) - count(m[a:])
	}

	if q.IsHead && q.Lines > 0 {
		i := a

		for n := 0; i < b && n < int(q.Lines); i++ {
			if m[i] == CR {
				n++
			}
		}

		b = min(i, b)
	}

	if q.IsTail && q.Lines > 0 {
		i, n := b-1, 0

		for ; i > a && n < int(q.Lines); i-- {
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

		// save last lines
		q.Values.Bytes = a
		q.Values.Lines = count(m) - n
	}

	return m[a:b]
}

func count(m mmap.MMap) int {
	v := bytes.Count(m, []byte{CR})

	if len(m) > 0 && m[len(m)-1] != CR {
		v++ // last line
	}

	return v
}
