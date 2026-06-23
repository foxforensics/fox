package unique

import (
	"github.com/zeebo/xxh3"
	"go.foxforensics.eu/fox/v4/internal/pkg/types"
)

type Hash struct {
	cache *types.Set
}

func ByHash() *Hash {
	return &Hash{
		cache: types.NewSet(),
	}
}

func (u *Hash) IsUnique(s string) bool {
	h := xxh3.HashString(s)

	if !u.cache.Has(h) {
		u.cache.Set(h)
		return true
	}

	return false
}
