package unique

import (
	"strings"
	"sync"

	"github.com/agnivade/levenshtein"
)

const window = 4096

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
	defer u.Unlock()

	// search latest strings first to improve performance
	for i := len(u.dedup) - 1; i >= max(0, len(u.dedup)-window); i-- {
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
