package text

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/cuhsat/fox/v4/internal/pkg/types/receipt"
)

var Banner = strings.TrimSpace(`
.-------.----.--.  .--.   .--. .--.--. .--.-. .--.-----.
|   ___/ .__. \  \/  /    |  |_|  |  | |  |  \|  |   _/
|   __|  |  |  >    <     |   _   |  | |  |   '  |  |
|  |   \ '--' /  /\  \    |  | |  |  '-'  |  |\  |  |
'--'    '----'--'  '--'   '--' '--'-------'--' '-'--'
`)

// standard output
var stdout io.Writer = os.Stdout

func Setup(w io.WriteCloser, err error) {
	stdout = w

	if err != nil {
		log.Fatalln(err)
	}
}

func Usage(msg string) error {
	_, _ = fmt.Println(Banner)
	_, _ = fmt.Println(msg)

	return nil
}

func Title(s string) {
	_, _ = fmt.Fprintln(stdout, bold.Sprint(s))
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
