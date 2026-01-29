package unique

import "github.com/agnivade/levenshtein"

type Unique struct {
	limit float64
	cache map[string]struct{}
}

func New(limit float64) *Unique {
	return &Unique{
		limit: limit,
		cache: make(map[string]struct{}),
	}
}

func (ls *Unique) IsUnique(s string) bool {
	for e := range ls.cache {
		d := levenshtein.ComputeDistance(e, s)

		// normalize distance
		if float64(d)/(float64(len(s)+len(e))) <= ls.limit {
			return false
		}
	}

	ls.cache[s] = struct{}{}
	return true
}
