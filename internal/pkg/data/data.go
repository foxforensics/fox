// Package data source:
// https://github.com/frizb/FirmwareReverseEngineering/blob/master/IdentifyingCompressionAlgorithms.md
// https://en.wikipedia.org/wiki/List_of_file_signatures
package data

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const tests = "../../../testdata/"

type Entry struct {
	Path string // Entry path
	Data []byte // Entry data
}

type Convert func([]byte) ([]byte, error)

type Deflate func([]byte) ([]byte, error)

type Extract func([]byte, string, string) []Entry

func HasMagic(b []byte, off int, m []byte) bool {
	if len(b) < off+len(m) {
		return false
	}

	return bytes.Equal(b[off:off+len(m)], m)
}

func Assert(b []byte) bool {
	return bytes.Equal(b, Fixture("fox.gs"))
}

func Fixture(name string) []byte {
	_, c, _, ok := runtime.Caller(0)

	if !ok {
		log.Fatalln("runtime error")
	}

	buf, err := os.ReadFile(filepath.Join(filepath.Dir(c), tests, name))

	if err != nil {
		log.Fatalln(err)
	}

	return buf
}
