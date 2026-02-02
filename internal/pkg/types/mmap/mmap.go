package mmap

import (
	"errors"
	"log"
	"os"
	"syscall"

	"github.com/cuhsat/go-mmap"
)

type MMap mmap.MMap

func Map(f *os.File) MMap {
	m, err := mmap.Map(f, mmap.RDONLY, 0)

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
