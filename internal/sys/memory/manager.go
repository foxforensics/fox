package memory

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
)

var mapped sync.Map

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
	slog.Debug(fmt.Sprintf("memory freed %s", k))

	if v, ok := mapped.LoadAndDelete(k); !ok {
		slog.Error(fmt.Sprintf("memory not found for %s", k))
	} else if m, ok := v.(MMap); ok {
		Unmap(m)
	}
}

func Purge() {
	slog.Debug("purging all memory")

	mapped.Range(func(k, v interface{}) bool {
		if m, ok := v.(MMap); !ok {
			Unmap(m)
			slog.Debug(fmt.Sprintf("memory purged for %s", k))
		}

		return true
	})
}
