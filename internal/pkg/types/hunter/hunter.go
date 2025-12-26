package hunter

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"maps"
	"regexp"
	"slices"
	"sync"

	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/evtx"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/journal"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heap"
)

const Level = 8

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

type Hunter struct {
	opts   *Options
	events chan *event.Event
}

func New(opts *Options) *Hunter {
	return &Hunter{
		opts:   opts,
		events: make(chan *event.Event, types.Size),
	}
}

func (htr *Hunter) Hunt(heaps <-chan *heap.Heap) <-chan *event.Event {
	if htr.opts.Sort {
		htr.events = htr.sort(htr.events)
	}

	go func(ch chan<- *event.Event) {
		for h := range heaps {
			if htr.opts.Verbose > 0 {
				log.Printf("hunt: parsing heap %s\n", h.String())
			}

			var wg sync.WaitGroup

			wg.Add(2)

			// hunt Windows Event Logs
			// hunt Linux Systemd journals
			go func(ch chan<- *event.Event) {
				defer wg.Done()
			}(ch)

			wg.Wait()

			h.Discard()
		}

		close(ch)
	}(htr.events)

	return htr.events
}

func (htr *Hunter) HuntCached(heaps <-chan *heap.Heap) <-chan *event.Event {
	return htr.events
}

func (htr *Hunter) huntEventLogs(h *heap.Heap) {
	r1 := bytes.NewReader(h.MMap())
	r2 := bytes.NewReader(h.MMap())

	for off := range htr.offset(r1, evtx.Regex) {
		for evt := range evtx.Carve(r2, off, htr.opts.Extensions) {
			htr.events <- evt
		}
	}
}

func (htr *Hunter) huntJournals(h *heap.Heap) {
	r := bytes.NewReader(h.MMap())

	for off := range htr.offset(r, journal.Regex) {
		for evt := range journal.Carve(h.MMap(), off, htr.opts.Extensions) {
			htr.events <- evt
		}
	}
}

func (htr *Hunter) sort(in <-chan *event.Event) chan *event.Event {
	out := make(chan *event.Event, cap(in))

	print("USE SORT")

	go func(out chan<- *event.Event) {
		cache := make(map[int64]*event.Event)

		for e := range in {
			cache[e.Time.UnixNano()] = e
		}

		print("CACHED")

		for _, k := range slices.Sorted(maps.Keys(cache)) {
			out <- cache[k]
		}

		print("SORTED")

		close(out)
	}(out)

	return out
}

func (htr *Hunter) offset(rs io.ReadSeeker, re *regexp.Regexp) <-chan int64 {
	out := make(chan int64, types.Size)

	go func(r *bufio.Reader) {
		var lst, off int64

		for loc := re.FindReaderIndex(r); loc != nil; loc = re.FindReaderIndex(r) {
			cur, _ := rs.Seek(0, io.SeekCurrent)
			off = lst + int64(loc[0])
			lst = cur - int64(r.Buffered())

			out <- off

			if htr.opts.Verbose > 2 {
				log.Printf("hunt: parsing offset 0x%08x\n", off)
			}
		}

		close(out)
	}(bufio.NewReader(rs))

	return out
}
