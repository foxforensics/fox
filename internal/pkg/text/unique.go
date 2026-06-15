package text

import (
	"strings"

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

type Distance struct {
	limit float64
	dedup []string
}

func ByDistance(limit float64) *Distance {
	return &Distance{
		limit: limit,
		dedup: make([]string, 0),
	}
}

func (u *Distance) IsUnique(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))

	// search latest strings first to improve performance
	for i := len(u.dedup) - 1; i >= 0; i-- {
		v := u.dedup[i]
		d := levenshtein.ComputeDistance(v, s)

		// normalize distance
		if float64(d)/(float64(len(s)+len(v))) <= u.limit {
			return false
		}
	}

	u.dedup = append(u.dedup, s)

	return true
}
