package data

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ulikunitz/xz"
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

	if !HasMagic(buf, 0, []byte{
		0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00,
	}) {
		return buf
	}

	r, err := xz.NewReader(bytes.NewReader(buf))

	if err != nil {
		log.Fatalln(err)
	}

	buf, err = io.ReadAll(r)

	if err != nil {
		log.Fatalln(err)
	}

	return buf
}
