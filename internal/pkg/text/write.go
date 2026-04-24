package text

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"go.foxforensics.dev/fox/v4/internal/pkg/file"
	"golang.org/x/term"

	"go.foxforensics.dev/fox/v4/internal"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/receipt"
)

var Banner = `
  .-------.----.--.  .--. 
  |   ___/ .__. \  \/  /   Created by Fox Forensics
  |   __|  |  |  >    <    https://foxforensics.eu
  |  |   \ '--' /  /\  \   Version %s
  '--'    '----'--'  '--'
`

// standard output
var stdout io.Writer = os.Stdout

func Setup(w io.WriteCloser, err error) {
	stdout = w

	if err != nil {
		log.Fatalln(err)
	}
}

func Usage(msg string) error {
	_, _ = fmt.Println(fmt.Sprintf(Banner, ver.Version))
	_, _ = fmt.Println(msg)

	return nil
}

func Title(s ...string) {
	w, _, err := term.GetSize(int(os.Stdin.Fd()))

	if err != nil {
		w = 78 // default
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

	_, _ = fmt.Fprint(stdout, Surface1.Sprint(" FOX "))
	_, _ = fmt.Fprint(stdout, Surface2.Sprintf(" %-*s ", w-29, title))
	_, _ = fmt.Fprint(stdout, Surface3.Sprintf(" %s ", stamp))
	_, _ = fmt.Fprintln(stdout)
}

func Match(s string, re *regexp.Regexp) {
	if re != nil && re.MatchString(s) {
		Write(MarkMatch(s, re))
	} else if re == nil {
		Write(s)
	}
}

func Write(f string, a ...any) {
	_, _ = fmt.Fprintf(stdout, fmt.Sprintf("%s\n", f), a...)
}

func Close(p string, r bool) {
	if v, is := stdout.(io.Closer); is {
		_ = v.Close()
	}

	if r && len(p) > 0 {
		err := receipt.Generate(p)

		if err != nil {
			log.Println(err)
		}
	}
}
