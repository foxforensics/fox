package storage

import "github.com/cuhsat/fox/v4/internal/pkg/types/event"

type Storage interface {
	String() string
	Write(*event.Event) error
	Close() error
}
