package hunter

import (
	"bytes"
	"log"
	"maps"
	"slices"

	"github.com/sourcegraph/conc"
	"github.com/sourcegraph/conc/pool"

	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/log/evtx"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/log/journal"
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
	Sort    bool
	Profile int
	Verbose int
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
	}

	return htr.events
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
	r := bytes.NewReader(h.Bytes())

	for off := range htr.findOffset(h.Bytes(), evtx.Chunk) {
		for evt := range evtx.Carve(r, off, cap(htr.events)) {
			htr.events <- evt
		}
	}
}

func (htr *Hunter) carveJournal(h *heap.Heap) {
	for off := range htr.findOffset(h.Bytes(), journal.Magic) {
		for evt := range journal.Carve(h.Bytes(), off, cap(htr.events)) {
			htr.events <- evt
		}
	}
}

func (htr *Hunter) findOffset(b, p []byte) <-chan int {
	out := make(chan int, 64*htr.opts.Profile)

	go func() {
		var off, idx int

		for {
			idx = bytes.Index(b[off:], p)

			if idx == -1 {
				break
			}

			out <- off + idx

			if htr.opts.Verbose > 2 {
				log.Printf("hunt: parsing offset 0x%08x\n", off+idx)
			}

			off += idx + len(p)
		}

		close(out)
	}()

	return out
}
