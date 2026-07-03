package memory

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
)

var (
	mapped sync.Map
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

func Alloc(k string, f *os.File) (MMap, error) {
	slog.Debug(fmt.Sprintf("memory alloc %s", k))

	m, err := Map(f)

	if err != nil {
		return m, err
	}

	mapped.Store(k, m)

	return m, nil
}

func Free(k string) {
	if v, ok := mapped.LoadAndDelete(k); !ok {
		slog.Error(fmt.Sprintf("memory not found for %s", k))
	} else if m, ok := v.(MMap); ok {
		if err := Unmap(m); err == nil {
			slog.Debug(fmt.Sprintf("memory freed %s", k))
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
				slog.Debug(fmt.Sprintf("memory purged for %s", k))
			} else {
				slog.Error(err.Error())
			}
		}

		return true
	})
}
