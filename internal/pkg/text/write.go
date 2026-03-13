package text

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/cuhsat/fox/v4/internal/pkg/types/receipt"
)

// standard output
var stdout io.Writer = os.Stdout

func Init(w io.Writer, err error) {
	stdout = w

	if err != nil {
		log.Fatalln(err)
	}
}

func Redirect(w io.Writer) {
	stdout = w
}

func Separator() string {
	return strings.Repeat("╌", width()-3)
}

func Framed(s string) {
	var sb strings.Builder

	l := strings.Repeat("─", width()-2)

	sb.WriteString(blue.Sprintf("┎%s┐\n", l))
	sb.WriteString(blue.Sprintf("┃ "))
	sb.WriteString(AsBold(s))
	sb.WriteString(strings.Repeat(" ", width()-len(s)-4))
	sb.WriteString(blue.Sprintf(" │\n"))
	sb.WriteString(blue.Sprintf("┠%s┘", l))

	_, _ = fmt.Fprintln(stdout, sb.String())
}

func Pretty(f string, a ...any) {
	_, _ = fmt.Fprintf(stdout, fmt.Sprintf("%s %s\n", Border, f), a...)
}

func Writeln(f string, a ...any) {
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

func width() int {
	n, _, err := term.GetSize(0)

	if err != nil {
		n = 78 // default term width
	}

	return n
}
