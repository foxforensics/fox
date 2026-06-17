package types

type Nil struct{}

type Set map[any]Nil

func (set *Set) Set(key any) {
	(*set)[key] = Nil{}
}

func (set *Set) Has(key any) bool {
	_, ok := (*set)[key]
	return ok
}
