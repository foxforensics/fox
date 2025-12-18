// Package writer provides a chain of custody file writer
package writer

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type Writer struct {
	sync.Mutex

	path string
	file *os.File
}

func New(path string) *Writer {
	var err error

	w := Writer{path: path}
	w.file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)

	if err != nil {
		log.Fatal(err)
	}

	return &w
}

func (w *Writer) Write(p []byte) (n int, err error) {
	w.Lock()
	defer w.Unlock()
	return w.file.Write(p)
}

func (w *Writer) Close() error {
	w.Lock()
	defer w.Unlock()
	_, err := w.file.Seek(0, io.SeekStart)

	if err != nil {
		return err
	}

	buf, err := io.ReadAll(w.file)

	if err != nil {
		return err
	}

	sum := fmt.Sprintf("%x", sha256.Sum256(buf))
	err = os.WriteFile(w.path+".sha256", []byte(sum), 0600)

	if err != nil {
		return err
	}

	return w.file.Close()
}
