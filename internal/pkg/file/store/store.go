package store

import "go.foxforensics.dev/fox/v4/internal/pkg/types/event"

type Store interface {
	String() string
	Store(*event.Event) error
	Close() error
}
