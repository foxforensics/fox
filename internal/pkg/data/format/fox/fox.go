package fox

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/cuhsat/fox/v4/internal/pkg/data/format/json"
)

const Head = "<FOX"
const Tail = "XOF>"

var cut = regexp.MustCompile(`^\s*[{}],?$`)
var rep = strings.NewReplacer(
	`  "`, `  `,
	`":`, `:`,
)

func Detect(b []byte) bool {
	s := string(b)
	return strings.HasPrefix(s, Head) && strings.HasSuffix(s, Tail)
}

func Format(b []byte) []byte {
	buf := bytes.NewBuffer(nil)
	res := json.Format(Unwrap(b))

	for i, line := range strings.Split(string(res), "\n") {
		if !cut.MatchString(line) {
			if i > 1 {
				buf.WriteRune('\n')
			}

			buf.WriteString(trim(line))
		}
	}

	return buf.Bytes()
}

func FromString(s string) string {
	return fmt.Sprintf("%s%s%s\n", Head, s, Tail)
}

func FromBytes(b []byte) []byte {
	return []byte(fmt.Sprintf("%s%s%s\n", Head, b, Tail))
}

func Unwrap(b []byte) []byte {
	s := string(b)
	s = strings.TrimPrefix(s, Head)
	s = strings.TrimSuffix(s, Tail)
	return []byte(s)
}

func trim(s string) string {
	s = strings.TrimSuffix(s[2:], " {")
	s = strings.TrimSuffix(s, "\":")
	s = strings.TrimPrefix(s, "\"")

	s = rep.Replace(s)

	return s
}
