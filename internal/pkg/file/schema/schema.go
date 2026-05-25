package schema

type Schema int

const (
	Raw Schema = iota
	Ecs
	Hec
)

func (shm Schema) String() string {
	return [...]string{"raw", "ecs", "hec"}[shm]
}
