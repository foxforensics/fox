package unique

import (
	"sync"

	"github.com/zeebo/xxh3"
	"go.foxforensics.eu/fox/v4/internal/pkg/types"
)

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
