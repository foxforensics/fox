package text

import (
	"fmt"
	"strings"

	"golang.org/x/term"
)

const (
	SP  = 0x20
	DEL = 0x7E
	LRE = 0x202A
	RLE = 0x202B
	PDF = 0x202C
	LRO = 0x202D
	RLO = 0x202E
	LRI = 0x2066
	RLI = 0x2067
	FSI = 0x2068
	PDI = 0x2069
)

func Line() string {
	return strings.Repeat("─", width())
}

func Title(line string) string {
	return Block(width(), line)
}

func Block(w int, lines ...string) string {
	var sb strings.Builder

	l := strings.Repeat("─", w-2)

	sb.WriteString(fmt.Sprintf("┌%s┐\n", l))

	for _, line := range lines {
		sb.WriteString(fmt.Sprintf("│ %-*s │\n", w-4, line))
	}

	sb.WriteString(fmt.Sprintf("└%s┘", l))

	return sb.String()
}

func ToAscii(s string) string {
	var sb strings.Builder

	for _, r := range s {
		if r < SP || r > DEL {
			sb.WriteRune('.')
		} else {
			sb.WriteRune(r)
		}
	}

	return sb.String()
}

func Sanitize(s string) string {
	var sb strings.Builder

	for _, r := range s {
		switch r { // mitigate CVE-2021-42574
		case LRE, RLE, LRO, RLO, LRI, RLI, FSI, PDF, PDI:
			sb.WriteRune('×')
		default:
			sb.WriteRune(r)
		}
	}

	return sb.String()
}

func Humanize(i int64) string {
	const m = int64(1024) // IEC prefix

	if i < m {
		return fmt.Sprintf("%db", i)
	}

	d, e := m, 0

	for n := i / m; n >= m; n /= m {
		d *= m
		e++
	}

	return fmt.Sprintf("%.1f%c", float64(i)/float64(d), "kmgtpezyrq"[e])
}

func width() int {
	w, _, err := term.GetSize(0)

	if err != nil {
		w = 78 // default term width
	}

	return w
}
