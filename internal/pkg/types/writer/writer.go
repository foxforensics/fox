// Package writer provides a chain of custody file writer
package writer

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"

	app "github.com/cuhsat/fox/v4/internal"
)

var header = strings.TrimSpace(`
FOX CHAIN OF CUSTODY RECEIPT %s
Time: %s
User: %s (%s)
Host: %s (%s)
Path: %s
Hash: %x SHA256
`)

type Writer struct {
	sync.Mutex

	path    string
	receipt bool
	file    *os.File
}

func New(path string, receipt bool) *Writer {
	var err error

	w := Writer{path: path, receipt: receipt}
	w.file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)

	if err != nil {
		log.Fatalln(err)
	}

	return &w
}

func (w *Writer) Close() error {
	w.Lock()
	defer w.Unlock()

	if w.receipt {
		err := w.WriteReceipt()

		if err != nil {
			log.Println(err)
		}
	}

	return w.file.Close()
}

func (w *Writer) Write(p []byte) (n int, err error) {
	w.Lock()
	defer w.Unlock()
	return w.file.Write(p)
}

func (w *Writer) WriteReceipt() error {
	hst, err := os.Hostname()

	if err != nil {
		return err
	}

	usr, err := user.Current()

	if err != nil {
		return err
	}

	abs, err := filepath.Abs(w.path)

	if err != nil {
		return err
	}

	_, err = w.file.Seek(0, io.SeekStart)

	if err != nil {
		return err
	}

	buf, err := io.ReadAll(w.file)

	if err != nil {
		return err
	}

	return os.WriteFile(w.path+".cc", []byte(fmt.Sprintf(header,
		app.Version[1:],
		time.Now().UTC(),
		usr.Name,
		usr.Username,
		hst,
		getMacAddr(),
		abs,
		sha256.Sum256(buf),
	)), 0600)
}

func getMacAddr() string {
	iff, err := net.Interfaces()

	if err == nil {
		for _, i := range iff {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				return i.HardwareAddr.String()
			}
		}
	}

	return ""
}
