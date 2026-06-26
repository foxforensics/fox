package parquet

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/parquet-go/parquet-go"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/event"
)

type Parquet struct {
	path   string
	file   *os.File
	writer *parquet.GenericWriter[event.Event]
}

func New(name string) (*Parquet, error) {
	var err error

	name = fmt.Sprintf("%s.parquet", name)

	prq := &Parquet{path: name}
	prq.file, err = os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		return nil, err
	}

	prq.writer = parquet.NewGenericWriter[event.Event](prq.file)

	return prq, nil
}

func (prq *Parquet) Run(ctx context.Context, ch <-chan *event.Event) {
	for {
		select {
		case <-ctx.Done():
			return

		case e := <-ch:
			if _, err := prq.writer.Write([]event.Event{*e}); err != nil {
				slog.Error(err.Error())
			}
		}
	}
}

func (prq *Parquet) String() string {
	return prq.path
}

func (prq *Parquet) Close() error {
	if prq.writer != nil {
		err := prq.writer.Close()

		if err != nil {
			return err
		}
	}

	if prq.file != nil {
		err := prq.file.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
