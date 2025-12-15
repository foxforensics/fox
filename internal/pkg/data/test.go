package data

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const Sample = "fox.txt"

func Assert(b []byte) bool {
	return bytes.Equal(b, Fixture(Sample))
}

func Fixture(name string) []byte {
	const dir = "../../../testdata"

	_, c, _, ok := runtime.Caller(0)

	if !ok {
		log.Fatalln("runtime error")
	}

	buf, err := os.ReadFile(filepath.Join(filepath.Dir(c), dir, name))

	if err != nil {
		log.Fatalln(err)
	}

	return buf
}
