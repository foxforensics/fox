package smap

import (
	"bufio"
	"bytes"
	"regexp"
	"runtime"
	"slices"

	"github.com/cuhsat/go-mmap"
	"github.com/sourcegraph/conc"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
	"github.com/cuhsat/fox/v4/internal/pkg/types/register"
)

const size = 1024 * 1024 * 4 // 4mb

var (
	sep = []byte("\n")
	tab = []byte{'\t'}
	exp = []byte("  ")
)

type action func(ch chan<- String, c *chunk)

type SMap []String

type String struct {
	Line  uint   // string line
	Group uint   // string group
	Bytes []byte // string data
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
			Line:  uint(len(s)) + 1,
			Bytes: bytes.Clone(b),
		})
	}

	return s
}

func (s SMap) Format() data.Format {
	if len(s) > 0 && len(register.Formats) > 0 {
		b := s[0].Bytes

		for _, f := range register.Formats {
			if f.Detect(b) {
				return f.Format
			}
		}
	}

	return nil
}

func (s SMap) Render() SMap {
	fn := s.Format() // check only first line

	return apply(func(ch chan<- String, c *chunk) {
		for _, s := range s[c.min:c.max] {
			if fn == nil {
				ch <- String{s.Line, s.Group, bytes.ReplaceAll(s.Bytes, tab, exp)}
				continue
			}

			for b := range bytes.SplitSeq(fn(s.Bytes), sep) {
				ch <- String{s.Line, s.Group, b}
			}
		}
	}, len(s))
}

func (s SMap) Grep(re *regexp.Regexp) SMap {
	return apply(func(ch chan<- String, c *chunk) {
		for _, str := range s[c.min:c.max] {
			if re.Match(str.Bytes) {
				ch <- str
			}
		}
	}, len(s))
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
		var wg conc.WaitGroup

		for _, c := range chunks(n) {
			wg.Go(func() {
				fn(ch, c)
			})
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
		if a.Group != b.Group {
			return int(a.Group - b.Group)
		}

		return int(a.Line - b.Line)
	})

	return s
}
