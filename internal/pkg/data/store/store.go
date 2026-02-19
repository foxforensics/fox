package store

import "github.com/cuhsat/fox/v4/internal/pkg/types/event"

type Store interface {
	String() string
	Store(*event.Event) error
	Close() error
}
