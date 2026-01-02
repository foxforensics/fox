package heap

import (
	"runtime"
	"sync"

	"github.com/edsrzf/mmap-go"

	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/smap"
)

type Context struct {
	Name   string
	Limit  *types.Limits
	Filter *types.Filters
}

type Heap struct {
	sync.RWMutex

	Name string // heap name
	Size int64  // heap size

	mmap mmap.MMap // memory map
	smap smap.SMap // string map
}

func New(ctx *Context, m mmap.MMap) *Heap {
	h := &Heap{
		Name: ctx.Name,
		Size: int64(len(m)),
		mmap: m,
	}

	// reduce mmap
	h.mmap = ctx.Limit.ReduceMMap(h.mmap)

	// reduce smap
	h.smap = ctx.Limit.ReduceSMap(smap.Map(h.mmap))

	// filter smap
	h.smap = ctx.Filter.FilterSMap(h.smap)

	return h
}

func (h *Heap) MMap() mmap.MMap {
	h.RLock()
	defer h.RUnlock()
	return h.mmap
}

func (h *Heap) SMap() smap.SMap {
	h.RLock()
	defer h.RUnlock()
	return h.smap
}

func (h *Heap) String() string {
	return h.Name
}

func (h *Heap) Discard() {
	h.Lock()

	// try to unmap original area
	_ = h.mmap.Unmap()

	h.Size = 0
	h.mmap = nil
	h.smap = nil

	h.Unlock()

	runtime.GC()
}
