package memory

import (
	"os"

	"go.foxforensics.eu/go-mmap"
)

type MMap mmap.MMap

func Map(f *os.File) (MMap, error) {
	m, err := mmap.Map(f, mmap.RDONLY, 0)
	return (MMap)(m), err
}

func Unmap(m MMap) error {
	return (*mmap.MMap)(&m).Unmap()
}
