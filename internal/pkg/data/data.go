// Package data source:
// https://github.com/frizb/FirmwareReverseEngineering/blob/master/IdentifyingCompressionAlgorithms.md
// https://en.wikipedia.org/wiki/List_of_file_signatures
package data

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const Sample = "fox.txt"
const Stream = ":"

const testdata = "../../../testdata"

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

func AddStream(p, s string) string {
	return fmt.Sprintf("%s%s%s", p, Stream, s)
}

func Assert(b []byte) bool {
	return bytes.Equal(b, Fixture(Sample))
}

func Fixture(name string) []byte {
	_, c, _, ok := runtime.Caller(0)

	if !ok {
		log.Fatalln("runtime error")
	}

	buf, err := os.ReadFile(filepath.Join(filepath.Dir(c), testdata, name))

	if err != nil {
		log.Fatalln(err)
	}

	return buf
}
