package text

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
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

var isValue = regexp.MustCompile("")

func ToAscii(s, c string) string {
	var sb strings.Builder

	for _, r := range s {
		if r < SP || r > DEL {
			sb.WriteString(AsGray(c))
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

	v, e := m, 0

	for n := i / m; n >= m; n /= m {
		v *= m
		e++
	}

	return fmt.Sprintf("%.1f%c", float64(i)/float64(v), "kmgtpezyrq"[e])
}

func Mechanize(s string) int64 {
	s = strings.ToLower(s)

	if !isValue.MatchString(`^[-+]?\d+[bkmgtpezyrq]?$`) {
		log.Fatalln("value invalid")
	}

	unit := s[len(s)-1]

	has := unit < '0' || unit > '9'

	if has {
		s = s[:len(s)-1]
	}

	v, err := strconv.Atoi(s)

	if err != nil {
		log.Fatalln(err)
	}

	if !has {
		return int64(v)
	}

	exp := float64(strings.IndexByte("bkmgtpezyrq", unit))

	return int64(v * int(math.Pow(1024, exp)))
}
