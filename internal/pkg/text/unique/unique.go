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

	_, ok := u.cache[h]

	if !ok {
		u.cache[h] = null{}
	}

	return !ok
}

func ByDistance(limit float64) *Distance {
	return &Distance{
		limit: limit,
		lines: make([]string, 0),
	}
}

func (u *Distance) IsUnique(s string) bool {
	for _, e := range u.lines {
		d := levenshtein.ComputeDistance(e, s)

		// normalize distance
		if float64(d)/(float64(len(s)+len(e))) <= u.limit {
			return false
		}
	}

	u.lines = append(u.lines, s)

	return true
}
