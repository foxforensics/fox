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

	"github.com/cuhsat/fox/v4/internal/pkg/files/format/evtx"
	"github.com/cuhsat/fox/v4/internal/pkg/files/format/journal"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heap"
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
	Extensions int
	Verbose    int
}

func Hunt(h *heap.Heap, opt *Options) <-chan *event.Event {
	ch := make(chan *event.Event, size)
	wg := sync.WaitGroup{}
	wg.Add(2)

	if opt.Verbose > 0 {
		log.Printf("hunt: parsing heap %s\n", h.String())
	}

	// hunt Windows Event Logs
	go func() {
		defer wg.Done()

		r1 := bytes.NewReader(h.MMap())
		r2 := bytes.NewReader(h.MMap())

		re := regexp.MustCompile(evtx.Chunk)

		for off := range offset(r1, re, opt) {
			for evt := range evtx.Decode(r2, off, opt.Extensions) {
				ch <- evt
			}
		}
	}()

	// hunt Linux Systemd journals
	go func() {
		defer wg.Done()

		r1 := bytes.NewReader(h.MMap())

		re := regexp.MustCompile(journal.Magic)

		for off := range offset(r1, re, opt) {
			for evt := range journal.Decode(h.MMap(), off, opt.Extensions) {
				ch <- evt
			}
		}
	}()

	// wait to close
	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func Sort(in <-chan *event.Event) chan *event.Event {
	out := make(chan *event.Event, cap(in))

	go func() {
		defer close(out)
		cache := make(map[int64]*event.Event)

		for e := range in {
			cache[e.Time.UnixNano()] = e
		}

		for _, k := range slices.Sorted(maps.Keys(cache)) {
			out <- cache[k]
		}
	}()

	return out
}

func offset(rs io.ReadSeeker, re *regexp.Regexp, opt *Options) <-chan int64 {
	ch := make(chan int64, size)

	go func(r *bufio.Reader) {
		var lst int64

		for loc := re.FindReaderIndex(r); loc != nil; loc = re.FindReaderIndex(r) {
			off, _ := rs.Seek(0, io.SeekCurrent)
			ch <- lst + int64(loc[0])
			lst = off - int64(r.Buffered())

			if opt.Verbose > 2 {
				log.Printf("hunt: parsing offset 0x%08x\n", loc[1])
			}
		}

		close(ch)
	}(bufio.NewReader(rs))

	return ch
}
