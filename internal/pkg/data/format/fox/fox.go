package fox

import (
	"bytes"
	"fmt"
	"strings"
)

var Head = []byte("<")
var Tail = []byte(">")

func Detect(b []byte) bool {
	return bytes.HasPrefix(b, Head) && bytes.HasSuffix(b, Tail)
}

func Format(b []byte) []byte {
	return Indent(b)
}

func FromString(s string) string {
	return fmt.Sprintf("%s%s%s\n", Head, s, Tail)
}

func FromBytes(b []byte) []byte {
	return []byte(fmt.Sprintf("%s%s%s\n", Head, b, Tail))
}

func Indent(b []byte) []byte {
	buf := bytes.NewBuffer(nil)

	b = bytes.TrimPrefix(b, Head)
	b = bytes.TrimSuffix(b, Tail)

	var numQuote int
	var numCurly int
	var numSquare int

	var lastSlash bool
	var lastBreak bool
	var nextCurly bool
	var nextSquare bool

	for i, c := range b {
		if i > 0 {
			lastSlash = b[i-1] == '\\'
			lastBreak = b[i-1] == '\n'
		}

		if i < len(b)-1 {
			nextCurly = b[i+1] == '}'
			nextSquare = b[i+1] == ']'
		}

		// literal char
		if numQuote%2 > 0 && (c != '"' || lastSlash) {
			buf.WriteByte(c)
			continue
		}

		switch c {
		case '{':
			numCurly += 1
			if numCurly+numSquare > 1 && !lastBreak && !nextCurly {
				buf.WriteByte('\n')
				buf.WriteString(tab(numCurly + numSquare))
			}
		case '}':
			numCurly -= 1
		case '[':
			numSquare += 1
			if numCurly+numSquare > 1 && !lastBreak && !nextSquare {
				buf.WriteByte('\n')
				buf.WriteString(tab(numCurly + numSquare))
			}
		case ']':
			numSquare -= 1
		case ':':
			buf.WriteByte(' ')
		case ',':
			if !lastBreak {
				buf.WriteByte('\n')
				buf.WriteString(tab(numCurly + numSquare))
			}
		case '"':
			numQuote += 1
		case ' ':
			// skip spaces
		default:
			buf.WriteByte(c)
		}
	}

	return buf.Bytes()
}

func tab(e int) string {
	return strings.Repeat("· ", e-1)
}
