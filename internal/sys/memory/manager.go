package memory

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
)

var (
	mapped  sync.Map
	counter atomic.Uint64
)

func Total() uint64 {
	n := uint64(0)

	mapped.Range(func(k, v any) bool {
		if m, ok := v.(MMap); ok {
			n += uint64(len(m))
		}
		return true
	})

	return n
}

func Alloc(f *os.File) (uint64, MMap, error) {
	m, err := Map(f)

	if err != nil {
		return 0, m, err
	}

	id := counter.Add(1)
	mapped.Store(id, m)

	slog.Debug(fmt.Sprintf("memory alloc for token %d", id))

	return id, m, nil
}

func Free(id uint64) {
	if v, ok := mapped.LoadAndDelete(id); !ok {
		slog.Error(fmt.Sprintf("memory not found for token %d", id))
	} else if m, ok := v.(MMap); ok {
		if err := Unmap(m); err == nil {
			slog.Debug(fmt.Sprintf("memory freed for token %d", id))
		} else {
			slog.Error(err.Error())
		}
	}
}

func Purge() {
	slog.Debug("purging all memory")

	mapped.Range(func(k, v any) bool {
		if m, ok := v.(MMap); ok {
			if err := Unmap(m); err == nil {
				slog.Debug(fmt.Sprintf("memory purged for %d", k))
			} else {
				slog.Error(err.Error())
			}
		}

		return true
	})
}
