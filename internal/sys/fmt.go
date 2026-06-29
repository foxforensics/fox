package sys

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2/v2"
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

var mask = regexp2.MustCompile(`^[-+]?\d+[bkmgtpezyr]?$`)

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

func Humanize(i uint64) string {
	const m = uint64(1024) // IEC prefix

	if i < m {
		return fmt.Sprintf("%db", i)
	}

	v, e := m, 0

	for n := i / m; n >= m; n /= m {
		v *= m
		e++
	}

	return fmt.Sprintf("%.1f%c", float64(i)/float64(v), "kmgtpezyr"[e])
}

func Mechanize(s string) (int64, bool) {
	s = strings.ToLower(s)

	if ok, _ := mask.MatchString(s); !ok {
		return 0, false
	}

	unit := s[len(s)-1]

	hasUnit := unit < '0' || unit > '9'

	if hasUnit {
		s = s[:len(s)-1] // cut unit
	}

	v, err := strconv.Atoi(s)

	if err != nil {
		return 0, false
	}

	if v < 0 {
		return 0, false
	}

	if !hasUnit {
		return int64(v), true
	}

	n := int64(v)

	for range strings.IndexByte("bkmgtpezyr", unit) {
		if n > math.MaxInt64/1024 || n < math.MinInt64/1024 {
			return 0, false // check for overflow and underflow
		}

		n *= 1024
	}

	return n, true
}
