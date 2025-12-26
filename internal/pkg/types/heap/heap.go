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
	Type   types.Heap
	Limit  *types.Limits
	Filter *types.Filters
}

type Heap struct {
	sync.RWMutex

	Name string     // heap name
	Type types.Heap // heap type
	Size int64      // heap size

	mmap mmap.MMap // memory map
	smap smap.SMap // string map
}

func New(ctx *Context, m mmap.MMap) *Heap {
	h := &Heap{
		Name: ctx.Name,
		Type: ctx.Type,
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
	switch h.Type {
	case types.Stdin:
		return "stdin"
	case types.Stdout:
		return "stdout"
	case types.Stderr:
		return "stderr"
	case types.Defined:
		return "input"
	default:
		return h.Name
	}
}

func (h *Heap) Discard() {
	h.Lock()

	if h.Type == types.Regular {
		_ = h.mmap.Unmap()
	}

	h.Size = 0
	h.mmap = nil
	h.smap = nil

	h.Unlock()

	runtime.GC()
}
