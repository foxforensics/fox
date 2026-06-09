package text

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dlclark/regexp2/v2"
	"golang.org/x/term"

	"go.foxforensics.dev/fox/v4/internal/pkg/file"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/receipt"
	"go.foxforensics.dev/fox/v4/internal/pkg/version"
)

var Banner = `
  .-------.----.--.  .--. 
  |   ___/ .__. \  \/  /   © %d by Fox Forensics
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
	_, err1 := fmt.Println(fmt.Sprintf(Banner, time.Now().Year(), version.Number))
	_, err2 := fmt.Println(msg)
	return errors.Join(err1, err2)
}

func SetOutput(wc io.WriteCloser, err error) {
	if err == nil {
		Stdout = &Writer{wc: wc}
	} else {
		log.Fatalln(err)
	}
}

func (w *Writer) Title(s ...string) {
	n, _, err := term.GetSize(int(os.Stdin.Fd()))

	if err != nil {
		n = 78 // default
	}

	title := s[0]
	title = strings.TrimPrefix(title, "/")
	title = strings.TrimSuffix(title, "/")
	title = strings.ReplaceAll(title, file.Separator, " ❱ ")
	title = strings.ReplaceAll(title, string(filepath.Separator), " ❱ ")

	if len(s) > 1 {
		title += " …"
	}

	stamp := time.Now().UTC().Format(time.RFC3339)

	w.Lock()
	_, _ = fmt.Fprint(w.wc, Surface1.Sprint(" FOX "))
	_, _ = fmt.Fprint(w.wc, Surface2.Sprintf(" %-*s ", n-29, title))
	_, _ = fmt.Fprint(w.wc, Surface3.Sprintf(" %s ", stamp))
	_, _ = fmt.Fprintln(w.wc)
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
			log.Println(err)
		}
	}
}
