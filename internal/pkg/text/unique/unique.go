package unique

import (
	"github.com/agnivade/levenshtein"
	"github.com/zeebo/xxh3"
)

type Unique interface {
	IsUnique(string) bool
}

type null struct{}

type Hash struct {
	cache map[uint64]null
}

type Distance struct {
	limit float64
	lines []string
}

func ByHash() *Hash {
	return &Hash{
		cache: make(map[uint64]null),
	}
}

func (u *Hash) IsUnique(s string) bool {
	h := xxh3.HashString(s)

	if _, ok := u.cache[h]; !ok {
		u.cache[h] = null{}
		return true
	}

	return false
}

func ByDistance(limit float64) *Distance {
	return &Distance{
		limit: limit,
		lines: make([]string, 0),
	}
}

func (u *Distance) IsUnique(s string) bool {
	for i := len(u.lines) - 1; i >= 0; i-- {
		l := u.lines[i]
		d := levenshtein.ComputeDistance(l, s)

		// normalize distance
		if float64(d)/(float64(len(s)+len(l))) <= u.limit {
			return false
		}
	}

	u.lines = append(u.lines, s)

	return true
}
