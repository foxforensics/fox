package stream

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

type Stream struct {
	path string

	f  *os.File
	ws []io.Writer
}

func New(path string, w io.Writer) *Stream {
	st := Stream{path: path}

	if len(path) > 0 {
		var err error

		st.f, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)

		if err != nil {
			log.Fatalln(err)
		}

		st.ws = append(st.ws, st.f)
	}

	if w != nil {
		st.ws = append(st.ws, w)
	}

	return &st
}

func (st *Stream) Write(p []byte) (n int, err error) {
	for _, w := range st.ws {
		_, err := w.Write(p)
		if err != nil {
			return 0, err
		}
	}

	return len(p), nil
}

func (st *Stream) Close() error {
	_, err := st.f.Seek(0, io.SeekStart)

	if err != nil {
		return err
	}

	buf, err := io.ReadAll(st.f)

	if err != nil {
		return err
	}

	sum := fmt.Sprintf("%x", sha256.Sum256(buf))
	err = os.WriteFile(st.path+".sha256", []byte(sum), 0600)

	if err != nil {
		return err
	}

	return st.f.Close()
}
