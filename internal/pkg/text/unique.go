package text

import (
	"strings"
	"sync"

	"github.com/agnivade/levenshtein"
	"github.com/zeebo/xxh3"

	"go.foxforensics.eu/fox/v4/internal/pkg/types"
)

type Unique interface {
	IsUnique(string) bool
}

type Hash struct {
	sync.Mutex
	cache types.Set
}

func ByHash() *Hash {
	return &Hash{
		cache: make(types.Set),
	}
}

func (u *Hash) IsUnique(s string) bool {
	h := xxh3.HashString(s)

	u.Lock()
	defer u.Unlock()

	if _, ok := u.cache[h]; !ok {
		u.cache.Set(h)
		return true
	}

	return false
}

type Distance struct {
	sync.Mutex
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

	u.Lock()

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

	u.Unlock()

	return true
}
