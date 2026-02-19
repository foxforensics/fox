package store

import "foxhunt.dev/fox/internal/pkg/types/event"

type Store interface {
	String() string
	Store(*event.Event) error
	Close() error
}
