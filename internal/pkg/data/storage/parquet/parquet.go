package parquet

import (
	"fmt"
	"log"
	"os"

	"github.com/cuhsat/fox/v4/internal/pkg/data/storage"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
	"github.com/parquet-go/parquet-go"
)

type File struct {
	path   string
	file   *os.File
	writer *parquet.GenericWriter[event.Event]
}

func New(name string) storage.Storage {
	var err error

	name = fmt.Sprintf("%s.parquet", name)

	f := &File{path: name}
	f.file, err = os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)

	if err != nil {
		log.Fatalln(err)
	}

	f.writer = parquet.NewGenericWriter[event.Event](f.file)

	return f
}

func (f *File) String() string {
	return f.path
}

func (f *File) Write(evt *event.Event) error {
	n, err := f.writer.Write([]event.Event{*evt})

	println(n)

	if err != nil {
		return err
	}

	return nil
}

func (f *File) Close() error {
	if f.file != nil {
		return f.file.Close()
	}

	return nil
}
