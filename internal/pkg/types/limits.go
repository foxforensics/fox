package types

import (
	"github.com/edsrzf/mmap-go"

	"github.com/cuhsat/fox/v4/internal/pkg/types/smap"
)

type Limits struct {
	IsHead bool // is head limit
	IsTail bool // is tail limit
	Lines  uint // lines count
	Bytes  uint // bytes count
}

func (l *Limits) ReduceMMap(m mmap.MMap) mmap.MMap {
	if l.IsHead && l.Bytes > 0 {
		r := make(mmap.MMap, min(l.Bytes, uint(len(m))))
		copy(r, m[:len(r)])
		return r
	}

	if l.IsTail && l.Bytes > 0 {
		r := make(mmap.MMap, min(uint(len(m)), l.Bytes))
		copy(r, m[max(len(m)-len(r), 0):])
		return r
	}

	return m
}

func (l *Limits) ReduceSMap(s smap.SMap) smap.SMap {
	if l.IsHead && l.Lines > 0 {
		r := make(smap.SMap, min(l.Lines, uint(len(s))))
		copy(r, s[:len(r)])
		return r
	}

	if l.IsTail && l.Lines > 0 {
		r := make(smap.SMap, min(uint(len(s)), l.Lines))
		copy(r, s[max(len(s)-len(r), 0):])
		return r
	}

	return s
}
