package smap

import (
	"bufio"
	"bytes"
	"regexp"
	"slices"

	"github.com/sourcegraph/conc"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
	"github.com/cuhsat/fox/v4/internal/pkg/types/register"
)

var Chunks = 2 // default

var sep = []byte("\n")
var tab = []byte{'\t'}
var exp = []byte("  ")

type action func(chan<- String, []String)

type SMap []String

type String struct {
	Line  uint   // string line
	Group uint   // string group
	Bytes []byte // string data
}

func Map(m []byte) (s SMap) {
	r := bufio.NewReader(bytes.NewReader(m))

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

	return
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

	return s.do(func(ch chan<- String, chk []String) {
		for _, str := range chk {
			if fn == nil {
				ch <- String{str.Line, str.Group, bytes.ReplaceAll(str.Bytes, tab, exp)}
				continue
			}

			for b := range bytes.SplitSeq(fn(str.Bytes), sep) {
				ch <- String{str.Line, str.Group, b}
			}
		}
	})
}

func (s SMap) Grep(re *regexp.Regexp) SMap {
	return s.do(func(ch chan<- String, chk []String) {
		for _, str := range chk {
			if re.Match(str.Bytes) {
				ch <- str
			}
		}
	})
}

func (s SMap) do(fn action) SMap {
	ch := make(chan String, len(s))

	go func() {
		var wg conc.WaitGroup

		for chk := range slices.Chunk(s, Chunks) {
			wg.Go(func() { fn(ch, chk) })
		}

		wg.Wait()

		close(ch)
	}()

	return sort(ch)
}

func sort(ch <-chan String) SMap {
	s := make(SMap, 0, len(ch))

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
