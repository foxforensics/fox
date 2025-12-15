package hunt

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"maps"
	"regexp"
	"slices"
	"sync"

	"github.com/cuhsat/fox/v4/internal/pkg/data/parser/evtx"
	"github.com/cuhsat/fox/v4/internal/pkg/data/parser/journal"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heapset"
)

const (
	Level = 8
	size  = 4096
)

var Paths = []string{
	"/Windows/System32/winevt/Logs",
	"/var/log/journal",
	"/run/log/journal",
}

type Options struct {
	Sort       bool
	Extensions int
	Verbose    int
}

func Hunt(hs *heapset.HeapSet, opt *Options) <-chan *event.Event {
	var wg sync.WaitGroup

	ch := make(chan *event.Event, size)

	if opt.Sort {
		ch = sort(ch)
	}

	wg.Add(2 * hs.Len())

	for _, h := range hs.Get() {
		if opt.Verbose > 0 {
			log.Printf("hunt: parsing heap %s\n", h.String())
		}

		r1 := bytes.NewReader(h.MMap())
		r2 := bytes.NewReader(h.MMap())
		r3 := bytes.NewReader(h.MMap())

		// hunt Windows Event Logs
		go func(ch chan<- *event.Event) {
			defer wg.Done()

			for off := range offset(r1, evtx.Regex, opt) {
				for evt := range evtx.Parse(r2, off, opt.Extensions) {
					ch <- evt
				}
			}
		}(ch)

		// hunt Linux Systemd journals
		go func(ch chan<- *event.Event) {
			defer wg.Done()

			for off := range offset(r3, journal.Regex, opt) {
				for evt := range journal.Parse(h.MMap(), off, opt.Extensions) {
					ch <- evt
				}
			}
		}(ch)
	}

	// wait to close
	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func sort(in <-chan *event.Event) chan *event.Event {
	out := make(chan *event.Event, cap(in))

	go func() {
		cache := make(map[int64]*event.Event)

		for e := range in {
			cache[e.Time.UnixNano()] = e
		}

		for _, k := range slices.Sorted(maps.Keys(cache)) {
			out <- cache[k]
		}

		close(out)
	}()

	return out
}

func offset(rs io.ReadSeeker, re *regexp.Regexp, opt *Options) <-chan int64 {
	out := make(chan int64, size)

	go func(r *bufio.Reader) {
		var lst, off int64

		for loc := re.FindReaderIndex(r); loc != nil; loc = re.FindReaderIndex(r) {
			cur, _ := rs.Seek(0, io.SeekCurrent)
			off = lst + int64(loc[0])
			lst = cur - int64(r.Buffered())

			out <- off

			if opt.Verbose > 2 {
				log.Printf("hunt: parsing offset 0x%08x\n", off)
			}
		}

		close(out)
	}(bufio.NewReader(rs))

	return out
}
