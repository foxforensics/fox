package text

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dlclark/regexp2/v2"

	"go.foxforensics.eu/fox/v4/internal/pkg/file"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/receipt"
	"go.foxforensics.eu/fox/v4/internal/pkg/version"
)

var Banner = `
  .-------.----.--.  .--. 
  |   ___/ .__. \  \/  /   © 2026 by Fox Forensics
  |   __|  |  |  >    <    https://foxforensics.eu
  |  |   \ '--' /  /\  \   Version %s
  '--'    '----'--'  '--'
`

// Stdout default output.
var Stdout = &Writer{wc: os.Stdout}

type Writer struct {
	sync.Mutex
	wc io.WriteCloser
}

func Usage(msg string) error {
	_, err1 := fmt.Println(fmt.Sprintf(Banner, version.Number))
	_, err2 := fmt.Println(msg)
	return errors.Join(err1, err2)
}

func SetOutput(wc io.WriteCloser, err error) {
	if err == nil {
		Stdout = &Writer{wc: wc}
	} else {
		slog.Error(err.Error())
	}
}

func (w *Writer) Header(s string) {
	s = strings.TrimPrefix(s, "/")
	s = strings.TrimSuffix(s, "/")
	s = strings.ReplaceAll(s, file.Separator, " > ")
	s = strings.ReplaceAll(s, string(filepath.Separator), " > ")

	w.Lock()
	_, _ = fmt.Fprintf(w.wc, "%s %s\n", Fox, AsBold(s))
	w.Unlock()
}

func (w *Writer) Match(s string, re *regexp2.Regexp) {
	if re != nil {
		if ok, _ := re.MatchString(s); !ok {
			return
		}
		s = MarkMatch(s, re)
	}

	w.Lock()
	_, _ = fmt.Fprintln(w.wc, s)
	w.Unlock()
}

func (w *Writer) Write(f string, a ...any) {
	w.Lock()
	_, _ = fmt.Fprintf(w.wc, fmt.Sprintf("%s\n", f), a...)
	w.Unlock()
}

func (w *Writer) Close(p string, r bool) {
	if v, is := w.wc.(io.Closer); is {
		_ = v.Close()
	}

	if r && len(p) > 0 {
		if err := receipt.Generate(p); err != nil {
			slog.Error(err.Error())
		}
	}
}
