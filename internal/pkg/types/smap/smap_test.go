package smap

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"

	"github.com/edsrzf/mmap-go"

	"github.com/cuhsat/fox/v4/internal/pkg/data/format/json"
	"github.com/cuhsat/fox/v4/internal/pkg/types/register"
)

func TestMain(m *testing.M) {
	register.Format("json", json.Detect, json.Format)

	os.Exit(m.Run())
}

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

func BenchmarkRender(b *testing.B) {
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

	for b.Loop() {
		s.Render()
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

func TestRender(t *testing.T) {
	f, m, err := fixture("format/fox.jsonl")

	if err != nil {
		t.Fatal(err)
	}

	s := Map(m).Render()

	if len(s) != 15 {
		t.Fatal("wrong length")
	}

	_ = m.Unmap()
	_ = f.Close()
}

func TestGrep(t *testing.T) {
	f, m, err := fixture("text/bible.txt")
	v := "Authorized King James Version\n"

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

	if s.String() != v {
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
