package schema

type Schema int

const (
	Raw Schema = iota
	Ecs
	Hec
)

func (shm Schema) String() string {
	switch shm {
	case Raw:
		return "raw"
	case Ecs:
		return "ecs"
	case Hec:
		return "hec"
	default:
		return "unknown"
	}
}
