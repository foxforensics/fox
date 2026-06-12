package hunter

import (
	"bytes"
	"errors"
	"io"
	"log"
	"maps"
	"slices"

	"github.com/sourcegraph/conc/pool"

	"go.foxforensics.eu/fox/v4/internal/pkg/file/binary/log/evtx"
	"go.foxforensics.eu/fox/v4/internal/pkg/file/binary/log/journal"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/event"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/heap"
)

var Block = 4096 * 16 // NTFS block size multiple
var Local = []string{
	"/Windows/System32/winevt/Logs",
	"/var/log/journal",
	"/run/log/journal",
}

type Options struct {
	Sort    bool
	Threads int
	Verbose int
}

type Hunter struct {
	opts   *Options
	cache  map[string]*event.Event
	events chan *event.Event
}

func New(opts *Options) *Hunter {
	return &Hunter{
		opts:   opts,
		cache:  make(map[string]*event.Event),
		events: make(chan *event.Event, opts.Threads*1024),
	}
}

func (htr *Hunter) Hunt(heaps <-chan *heap.Heap) <-chan *event.Event {
	go func() {
		defer close(htr.events)

		p := pool.New().WithMaxGoroutines(htr.opts.Threads)

		for h := range heaps {
			p.Go(func() {
				if htr.opts.Verbose > 0 {
					log.Printf("hunt: carving heap %s\n", h.String())
				}

				htr.carve(h)
			})
		}

		p.Wait()
	}()

	if htr.opts.Sort {
		return htr.sort()
	}

	return htr.events
}

func (htr *Hunter) sort() <-chan *event.Event {
	sorted := make(chan *event.Event, cap(htr.events))

	go func() {
		defer close(sorted)

		for e := range htr.events {
			htr.cache[e.SortKey()] = e
		}

		for _, k := range slices.Sorted(maps.Keys(htr.cache)) {
			sorted <- htr.cache[k]
		}
	}()

	return sorted
}

func (htr *Hunter) carve(h *heap.Heap) {
	defer h.Discard()

	p := pool.New().WithMaxGoroutines(htr.opts.Threads)

	p.Go(func() {
		htr.carveEvtx(h)
	})

	p.Go(func() {
		htr.carveJournal(h)
	})

	p.Wait()
}

func (htr *Hunter) carveEvtx(h *heap.Heap) {
	sr := io.NewSectionReader(h.Reader(), 0, int64(h.Size))

	for off := range htr.findOffset(h, evtx.Chunk) {
		for evt := range evtx.Carve(sr, off, cap(htr.events)) {
			htr.events <- evt
		}
	}
}

func (htr *Hunter) carveJournal(h *heap.Heap) {
	sr := io.NewSectionReader(h.Reader(), 0, int64(h.Size))

	for off := range htr.findOffset(h, journal.Magic) {
		for evt := range journal.Carve(sr, off, cap(htr.events)) {
			htr.events <- evt
		}
	}
}

func (htr *Hunter) findOffset(h *heap.Heap, seq []byte) <-chan int64 {
	out := make(chan int64, 64*htr.opts.Threads)

	go func(r io.ReaderAt, n uint64) {
		var off, idx int64

		blk := make([]byte, Block)

		for off < int64(n) {
			n, err := r.ReadAt(blk, off)

			if errors.Is(err, io.EOF) {
				break
			}

			if err != nil {
				log.Println(err)
			}

			idx = int64(bytes.Index(blk, seq))

			if idx == -1 {
				off += int64(n)
				continue
			}

			off += idx

			out <- off

			if htr.opts.Verbose > 2 {
				log.Printf("hunt: found at offset 0x%08x\n", off)
			}

			off += int64(len(seq))
		}

		close(out)
	}(h.Reader(), h.Size)

	return out
}
