package unique

import (
	"strings"
	"sync"

	"github.com/agnivade/levenshtein"
)

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
