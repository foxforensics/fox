package parquet

import (
	"fmt"
	"log"
	"os"

	"github.com/cuhsat/fox/v4/internal/pkg/data/store"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
	"github.com/parquet-go/parquet-go"
)

type File struct {
	path string
	f    *os.File
	w    *parquet.GenericWriter[event.Event]
}

func New(name string) store.Store {
	var err error

	name = fmt.Sprintf("%s.parquet", name)

	f := &File{path: name}
	f.f, err = os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		log.Fatalln(err)
	}

	f.w = parquet.NewGenericWriter[event.Event](f.f)

	return f
}

func (f *File) String() string {
	return f.path
}

func (f *File) Store(evt *event.Event) error {
	_, err := f.w.Write([]event.Event{*evt})

	return err
}

func (f *File) Close() error {
	if f.w != nil {
		err := f.w.Close()

		if err != nil {
			return err
		}
	}

	if f.f != nil {
		err := f.f.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
