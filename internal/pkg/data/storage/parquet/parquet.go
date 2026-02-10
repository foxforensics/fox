package parquet

import (
	"github.com/cuhsat/fox/v4/internal/pkg/data/storage"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
)

type File struct {
	path string
}

func New(path string) storage.Storage {
	f := &File{path: path}
	return f
}

func (f *File) String() string {
	return f.path + ".parquet"
}

func (f *File) Write(evt *event.Event) error {
	return nil
}

func (f *File) Close() error {
	return nil
}
