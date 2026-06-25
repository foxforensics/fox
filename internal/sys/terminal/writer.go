package terminal

import (
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dlclark/regexp2/v2"
	"go.foxforensics.eu/fox/v4/internal/sys/receipt"
)

type Writer struct {
	sync.Mutex
	wc io.WriteCloser
}

func NewWriter(wc io.WriteCloser) *Writer {
	return &Writer{wc: wc}
}

func (w *Writer) Header(s string) {
	s = strings.TrimPrefix(s, "/")
	s = strings.TrimSuffix(s, "/")
	s = strings.ReplaceAll(s, ":", " > ")
	s = strings.ReplaceAll(s, string(filepath.Separator), " > ")

	w.Lock()
	_, err := fmt.Fprintf(w.wc, "%s %s\n", Fox, AsBold(s))
	w.Unlock()

	if err != nil {
		slog.Error(err.Error())
	}
}

func (w *Writer) Match(s string, re *regexp2.Regexp) {
	if re != nil {
		if ok, _ := re.MatchString(s); !ok {
			return
		}
		s = MarkMatch(s, re)
	}

	w.Lock()
	_, err := fmt.Fprintln(w.wc, s)
	w.Unlock()

	if err != nil {
		slog.Error(err.Error())
	}
}

func (w *Writer) Write(f string, a ...any) {
	w.Lock()
	_, err := fmt.Fprintf(w.wc, f+"\n", a...)
	w.Unlock()

	if err != nil {
		slog.Error(err.Error())
	}
}

func (w *Writer) Close(p string, r bool) {
	if v, is := w.wc.(io.Closer); is {
		if err := v.Close(); err != nil {
			slog.Error(err.Error())
		}
	}

	if r && len(p) > 0 {
		if err := receipt.Generate(p); err != nil {
			slog.Error(err.Error())
		}
	}
}
