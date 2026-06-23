package types

import "sync"

type Nil struct{}

type Set struct {
	sync.RWMutex
	m map[any]Nil
}

func NewSet() *Set {
	return &Set{m: make(map[any]Nil)}
}

func (set *Set) Len() int {
	set.RLock()
	defer set.RUnlock()
	return len((*set).m)
}

func (set *Set) Set(key any) {
	set.Lock()
	(*set).m[key] = Nil{}
	set.Unlock()
}

func (set *Set) Has(key any) bool {
	set.RLock()
	_, ok := (*set).m[key]
	set.RUnlock()
	return ok
}
