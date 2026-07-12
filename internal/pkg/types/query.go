package types

import (
	"bytes"
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2/v2"
	"go.foxforensics.eu/fox/v5/internal/pkg/smap"
)

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

	if ok, err := limit.MatchString(s); !ok {
		if err != nil {
			slog.Debug(err.Error())
		}
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

func (q *Query) Reduce(b []byte) ([]byte, bool) {
	var x, y = 0, len(b)

	if !q.IsHead && !q.IsTail {
		return b, false
	}

	if q.IsHead && q.Bytes > 0 {
		y = min(int(q.Bytes), y)
	}

	if q.IsTail && q.Bytes > 0 {
		x = max(len(b)-int(q.Bytes), 0)

		// save last lines
		q.Values.Bytes = x
		q.Values.Lines = count(b) - count(b[x:])
	}

	if q.IsHead && q.Lines > 0 {
		i := x

		for n := 0; i < y && n < int(q.Lines); i++ {
			if b[i] == '\n' {
				n++
			}
		}

		y = min(i, y)
	}

	if q.IsTail && q.Lines > 0 {
		i, n := y-1, 0

		for ; i > x && n < int(q.Lines); i-- {
			if b[i-1] == '\n' {
				n++
			}
		}

		x = max(i, x)

		if x > 0 {
			x++ // skip linebreak
		}

		if i == 0 {
			n++ // add first line
		}

		// save last lines
		q.Values.Bytes = x
		q.Values.Lines = count(b) - n
	}

	return b[x:y], true
}

func count(b []byte) int {
	v := bytes.Count(b, []byte{'\n'})

	if len(b) > 0 && b[len(b)-1] != '\n' {
		v++ // last line
	}

	return v
}
