package mmap

import (
	"errors"
	"io"
	"log"
	"os"
	"syscall"

	"github.com/edsrzf/mmap-go"
)

type MMap mmap.MMap

func Map(f *os.File) MMap {
	m, err := mmap.Map(f, mmap.RDONLY, 0)

	if err != nil {
		log.Fatalln(err)
	}

	return (MMap)(m)
}

func Remap(r io.ReaderAt, size int) MMap {
	m, err := mmap.MapRegion(nil, size, mmap.RDWR, mmap.ANON, 0)

	if err != nil {
		log.Fatalln(err)
	}

	_, err = r.ReadAt(m, 0)

	if err != nil {
		log.Fatalln(err)
	}

	return (MMap)(m)
}

func Unmap(m MMap) {
	if err := (*mmap.MMap)(&m).Unmap(); err != nil {
		if !errors.Is(err, syscall.EINVAL) {
			log.Println(err)
		}
	}
}
