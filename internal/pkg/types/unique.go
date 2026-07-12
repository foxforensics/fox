package types

import (
	"sync"

	"github.com/zeebo/xxh3"
)

type Nil struct{}

type Unique struct {
	m sync.Map
}

func NewUnique() *Unique {
	return new(Unique)
}

func (u *Unique) Is(s string) bool {
	_, ok := u.m.LoadOrStore(xxh3.HashString(s), Nil{})
	return !ok
}
