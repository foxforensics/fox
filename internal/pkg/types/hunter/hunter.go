package hunter

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"maps"
	"regexp"
	"slices"

	"github.com/sourcegraph/conc"
	"github.com/sourcegraph/conc/pool"

	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/evtx"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/journal"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heap"
)

const Database = "fox.db"

var Local = []string{
	"/Windows/System32/winevt/Logs",
	"/var/log/journal",
	"/run/log/journal",
}

type Options struct {
	Sort       bool
	Extensions int
	Profile    int
	Verbose    int
}

type Hunter struct {
	opts   *Options
	cache  map[uint64]*event.Event
	events chan *event.Event
}

func New(opts *Options) *Hunter {
	return &Hunter{
		opts:   opts,
		cache:  make(map[uint64]*event.Event),
		events: make(chan *event.Event, opts.Profile*1024),
	}
}

func (htr *Hunter) Hunt(heaps <-chan *heap.Heap) <-chan *event.Event {
	go func() {
		p := pool.New().WithMaxGoroutines(htr.opts.Profile)

		for h := range heaps {
			p.Go(func() {
				if htr.opts.Verbose > 0 {
					log.Printf("hunt: carving heap %s\n", h.String())
				}

				htr.carve(h)
			})
		}

		p.Wait()

		close(htr.events)
	}()

	if htr.opts.Sort {
		return htr.sort()
	} else {
		return htr.events
	}
}

func (htr *Hunter) sort() <-chan *event.Event {
	sorted := make(chan *event.Event, cap(htr.events))

	go func() {
		for e := range htr.events {
			htr.cache[e.Hash()] = e
		}

		for _, k := range slices.Sorted(maps.Keys(htr.cache)) {
			sorted <- htr.cache[k]
		}

		close(sorted)
	}()

	return sorted
}

func (htr *Hunter) carve(h *heap.Heap) {
	defer h.Discard()

	var wg conc.WaitGroup

	wg.Go(func() {
		htr.carveEvtx(h)
	})

	wg.Go(func() {
		htr.carveJournal(h)
	})

	wg.Wait()
}

func (htr *Hunter) carveEvtx(h *heap.Heap) {
	r1 := bytes.NewReader(h.MMap())
	r2 := bytes.NewReader(h.MMap())

	for off := range htr.findOffset(r1, evtx.Regex) {
		for evt := range evtx.Carve(r2, off, htr.opts.Extensions, cap(htr.events)) {
			htr.events <- evt
		}
	}
}

func (htr *Hunter) carveJournal(h *heap.Heap) {
	r := bytes.NewReader(h.MMap())

	for off := range htr.findOffset(r, journal.Regex) {
		for evt := range journal.Carve(h.MMap(), off, htr.opts.Extensions, cap(htr.events)) {
			htr.events <- evt
		}
	}
}

func (htr *Hunter) findOffset(rs io.ReadSeeker, re *regexp.Regexp) <-chan int64 {
	out := make(chan int64, 64*htr.opts.Profile)

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
