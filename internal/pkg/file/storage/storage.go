package storage

import "go.foxforensics.dev/fox/v4/internal/pkg/types/event"

type Storage interface {
	Store(*event.Event) error
	String() string
	Close() error
}
