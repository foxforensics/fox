package smap

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"

	"github.com/cuhsat/go-mmap"
)

func BenchmarkMap(b *testing.B) {
	f, m, err := fixture("text/bible.txt")

	if err != nil {
		b.Fatal(err)
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	defer func(m mmap.MMap) {
		_ = m.Unmap()
	}(m)

	for b.Loop() {
		Map(m)
	}
}

func BenchmarkGrep(b *testing.B) {
	f, m, err := fixture("text/bible.txt")

	if err != nil {
		b.Fatal(err)
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	defer func(m mmap.MMap) {
		_ = m.Unmap()
	}(m)

	s := Map(m)

	re := regexp.MustCompile(".*")

	for b.Loop() {
		s.Grep(re)
	}
}

func TestMap(t *testing.T) {
	f, m, err := fixture("text/bible.txt")

	if err != nil {
		t.Fatal(err)
	}

	if len(Map(m)) != 31107 {
		t.Fatal("wrong size")
	}

	_ = m.Unmap()
	_ = f.Close()
}

func TestGrep(t *testing.T) {
	f, m, err := fixture("text/bible.txt")
	v := "Authorized King James Version"

	if err != nil {
		t.Fatal(err)
	}

	re := regexp.MustCompile("King James")

	s := Map(m).Grep(re)

	_ = m.Unmap()
	_ = f.Close()

	if len(s) != 1 {
		t.Fatal("wrong length")
	}

	if string(s[0].Bytes) != v {
		t.Fatal("wrong string")
	}
}

func fixture(name string) (*os.File, mmap.MMap, error) {
	_, c, _, ok := runtime.Caller(0)

	if !ok {
		return nil, nil, errors.New("error")
	}

	p := filepath.Join(filepath.Dir(c), "..", "..", "..", "..", "testdata", name)

	f, err := os.OpenFile(p, os.O_RDONLY, 0400)

	if err != nil {
		return nil, nil, err
	}

	m, err := mmap.Map(f, mmap.RDONLY, 0)

	if err != nil {
		return nil, nil, err
	}

	return f, m, nil
}
