package types

import (
	"strings"
	"sync"

	"github.com/agnivade/levenshtein"
	"github.com/zeebo/xxh3"
)

type Mode int

const (
	Hash Mode = iota
	Distance
)

const window = 4096

type Unique struct {
	sync.Mutex

	mode  Mode
	limit float64
	dedup []string
	cache sync.Map
}

func NewUnique(mode Mode) *Unique {
	return &Unique{
		mode:  mode,
		dedup: make([]string, 0),
	}
}

func (u *Unique) SetLimit(v float64) {
	u.limit = v
}

func (u *Unique) IsUnique(s string) bool {
	switch u.mode {
	case Hash:
		return u.byHash(s)
	case Distance:
		return u.byDistance(s)
	}

	return false
}

func (u *Unique) byHash(s string) bool {
	_, ok := u.cache.LoadOrStore(xxh3.HashString(s), Nil{})
	return !ok
}

func (u *Unique) byDistance(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))

	u.Lock()
	defer u.Unlock()

	// search latest strings first to improve performance
	for i := len(u.dedup) - 1; i >= max(0, len(u.dedup)-window); i-- {
		v := u.dedup[i]
		d := levenshtein.ComputeDistance(v, s)

		// normalize distance
		if float64(d)/(float64(max(len(s), len(v)))) <= u.limit {
			return false
		}
	}

	u.dedup = append(u.dedup, s)

	return true
}
