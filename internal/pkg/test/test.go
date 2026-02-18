package test

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/klauspost/compress/zstd"
)

const Sample = "fox.txt"

const root = "../../../testdata"

func Assert(b []byte) bool {
	return bytes.Equal(b, Fixture(Sample))
}

func FixtureDeflate(name string) string {
	b, err := os.ReadFile(FixtureFile(name))

	if err != nil {
		log.Fatalln(err)
	}

	t, err := os.CreateTemp("", "fox-*")

	if err != nil {
		log.Fatalln(err)
	}

	_, err = t.Write(deflate(b))

	if err != nil {
		log.Fatalln(err)
	}

	_ = t.Close()

	return t.Name()
}

func FixtureFile(name string) string {
	_, c, _, ok := runtime.Caller(0)

	if !ok {
		log.Fatalln("runtime error")
	}

	return filepath.Join(filepath.Dir(c), root, name)
}

func FixtureDir(names []string) []string {
	var v []string

	for _, name := range names {
		v = append(v, FixtureFile(name))
	}

	return v
}

func FixtureRaw(name string) []byte {
	buf, err := os.ReadFile(FixtureFile(name))

	if err != nil {
		log.Fatalln(err)
	}

	return buf
}

func Fixture(name string) []byte {
	return deflate(FixtureRaw(name))
}

func deflate(b []byte) []byte {
	if !bytes.Equal(b[:4], []byte{0x28, 0xB5, 0x2F, 0xFD}) {
		return b
	}

	r, err := zstd.NewReader(bytes.NewReader(b))

	if err != nil {
		log.Fatalln(err)
	}

	defer r.Close()

	b, err = io.ReadAll(r)

	if err != nil {
		log.Fatalln(err)
	}

	return b
}
