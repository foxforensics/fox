package storage

import (
	"go.foxforensics.eu/fox/v4/internal/cmd/hunt/event"
)

type Storage interface {
	Store(*event.Event) error
	String() string
	Close() error
}
