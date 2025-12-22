package fson

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/cuhsat/fox/v4/internal/pkg/types/buffer"
)

func Detect(b []byte) bool {
	return b[0] == '<' && b[len(b)-1] == '>'
}

func Format(b []byte) []byte {
	buf := bytes.NewBuffer(nil)

	// trim head and tail marks
	b = b[1 : len(b)-2]

	var numQuote int
	var numLevel int

	var prevSlash bool
	var prevBreak bool
	var nextBrace bool

	for i, c := range b {
		// look back
		if i > 0 {
			prevSlash = b[i-1] == '\\'
			prevBreak = b[i-1] == '\n'
		}

		// look ahead
		if i < len(b)-1 {
			nextBrace = b[i+1] == '}' || b[i+1] == ']'
		}

		// literal char
		if numQuote%2 > 0 && (c != '"' || prevSlash) {
			buf.WriteByte(c)
			continue
		}

		switch c {
		case ' ': // skip spaces
		case '{', '[':
			numLevel += 1
			if numLevel > 1 && !prevBreak && !nextBrace {
				indent(buf, numLevel-1)
			}
		case '}', ']':
			numLevel -= 1
		case '"':
			numQuote += 1
		case ',':
			indent(buf, numLevel-1)
		case ':':
			buf.WriteByte(' ')
		default:
			buf.WriteByte(c)
		}
	}

	buf.WriteByte('\n')
	buf.WriteString(buffer.Sep)

	return buf.Bytes()
}

func FromString(s string) string {
	return fmt.Sprintf("<%s>\n", s)
}

func FromBytes(b []byte) []byte {
	return []byte(FromString(string(b)))
}

func indent(b *bytes.Buffer, e int) {
	b.WriteByte('\n')
	b.WriteString(strings.Repeat("· ", e))
}
