package smap

import (
	"bufio"
	"bytes"
	"cmp"
	"slices"

	"github.com/dlclark/regexp2/v2"
	"github.com/sourcegraph/conc/iter"
)

const Size = 1024 * 1024 * 4 // 4m

var Threads = 2 // default

var tab = []byte{'\t'}
var spc = []byte("  ")

type action func(chan<- String, *String)

type SMap []String

type String struct {
	Line  uint   // string line
	Group uint   // string group
	Bytes []byte // string data
}

func Map(m []byte) (s SMap) {
	r := bufio.NewReaderSize(bytes.NewReader(m), Size)

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

func (s SMap) Render() SMap {
	return s.do(func(ch chan<- String, str *String) {
		ch <- String{str.Line, str.Group, bytes.ReplaceAll(str.Bytes, tab, spc)}
	})
}

func (s SMap) Grep(re *regexp2.Regexp) SMap {
	return s.do(func(ch chan<- String, str *String) {
		if ok, _ := re.MatchString(string(str.Bytes)); ok {
			ch <- *str
		}
	})
}

func (s SMap) do(fn action) SMap {
	ch := make(chan String, max(1, Threads*64))

	go func(chan<- String) {
		it := iter.Iterator[String]{
			MaxGoroutines: Threads,
		}

		it.ForEach(s, func(s *String) {
			fn(ch, s)
		})

		close(ch)
	}(ch)

	return sort(ch)
}

func sort(ch <-chan String) SMap {
	s := make(SMap, 0, len(ch))

	for str := range ch {
		s = append(s, str)
	}

	slices.SortStableFunc(s, func(a, b String) int {
		if a.Group != b.Group {
			return cmp.Compare(a.Group, b.Group)
		}

		return cmp.Compare(a.Line, b.Line)
	})

	return s
}
