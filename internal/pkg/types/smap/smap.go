package smap

import (
	"bufio"
	"bytes"
	"encoding/json"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"sync"

	"github.com/edsrzf/mmap-go"
)

const size = 1024 * 1024 * 4 // 4mb

type action func(ch chan<- String, c *chunk)

type SMap []String

type String struct {
	Nr  uint   // string nr
	Grp uint   // string group
	Str string // string data
}

type chunk struct {
	min uint // chunk start
	max uint // chunk end
}

func Map(m mmap.MMap) SMap {
	s := make(SMap, 0)

	r := bufio.NewReaderSize(bytes.NewReader(m), size)

	for {
		b, _, err := r.ReadLine()

		if err != nil {
			break
		}

		s = append(s, String{
			Nr:  uint(len(s)) + 1,
			Str: string(b),
		})
	}

	return s
}

func (s SMap) String() string {
	var sb strings.Builder

	for _, str := range s {
		sb.WriteString(str.Str)
		sb.WriteRune('\n')
	}

	return sb.String()
}

func (s SMap) Render(e int) SMap {
	tab := strings.Repeat(" ", e)
	return apply(func(ch chan<- String, c *chunk) {
		for _, s := range s[c.min:c.max] {
			ch <- String{s.Nr, s.Grp, expand(s.Str, tab)}
		}
	}, len(s))
}

func (s SMap) Format(e int) SMap {
	tab := strings.Repeat(" ", e)
	return apply(func(ch chan<- String, c *chunk) {
		var buf bytes.Buffer

		for _, s := range s[c.min:c.max] {
			buf.Reset()

			if json.Indent(&buf, []byte(s.Str), "", tab) != nil {
				ch <- String{s.Nr, s.Grp, s.Str}
				continue
			}

			for b := range bytes.SplitSeq(buf.Bytes(), []byte("\n")) {
				ch <- String{s.Nr, s.Grp, string(b)}
			}
		}
	}, len(s))
}

func (s SMap) Grep(re *regexp.Regexp) SMap {
	return apply(func(ch chan<- String, c *chunk) {
		for _, s := range s[c.min:c.max] {
			if re.MatchString(s.Str) {
				ch <- s
			}
		}
	}, len(s))
}

func (s SMap) CanFormat() bool {
	if len(s) == 0 {
		return false
	}

	return json.Valid([]byte(s[0].Str))
}

func chunks(n int) (c []*chunk) {
	m := min(runtime.NumCPU(), n)

	for i := range m {
		c = append(c, &chunk{
			min: uint(i * n / m),
			max: uint(((i + 1) * n) / m),
		})
	}

	return
}

func apply(fn action, n int) SMap {
	ch := make(chan String, n)

	go func() {
		var wg sync.WaitGroup

		for _, c := range chunks(n) {
			wg.Add(1)

			go func() {
				fn(ch, c)
				wg.Done()
			}()
		}

		wg.Wait()

		close(ch)
	}()

	return sort(ch)
}

func sort(ch <-chan String) SMap {
	s := make(SMap, 0)

	for str := range ch {
		s = append(s, str)
	}

	slices.SortStableFunc(s, func(a, b String) int {
		if a.Grp != b.Grp {
			return int(a.Grp - b.Grp)
		} else {
			return int(a.Nr - b.Nr)
		}
	})

	return s
}

func expand(s, t string) string {
	return strings.ReplaceAll(s, "\t", t)
}
