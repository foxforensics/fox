package mmap

import (
	"errors"
	"log/slog"
	"os"
	"syscall"

	"go.foxforensics.eu/go-mmap"
)

type MMap mmap.MMap

func Map(f *os.File) (MMap, error) {
	m, err := mmap.Map(f, mmap.RDONLY, 0)
	return (MMap)(m), err
}

func Unmap(m MMap) {
	if err := (*mmap.MMap)(&m).Unmap(); err != nil {
		if !errors.Is(err, syscall.EINVAL) {
			slog.Error(err.Error())
		}
	}
}
