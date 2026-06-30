package mft

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"time"

	"go.foxforensics.eu/fox/v4/internal/lib"
	"go.foxforensics.eu/go-mft"
)

const cluster = 4096 // sane size
const timeout = 60

func Detect(b []byte) bool {
	return lib.HasMagic(b, 0, []byte{
		'F', 'I', 'L', 'E',
	})
}

func Convert(b []byte) ([]byte, error) {
	dt, err := mft.BuildDirectoryTree(bytes.NewReader(b), "c")

	if err != nil {
		return b, err
	}

	ch := make(chan mft.UsefulMftFields)

	go mft.ParseMftRecords(bytes.NewReader(b), cluster, dt, &ch)

	v := make([]mft.UsefulMftFields, 0, len(b)/1024)

	for {
		select {
		// this only a fallback if something goes wrong within the parser
		case <-time.After(time.Second * timeout):
			slog.Error("mft: parsing timed out")
			close(ch) // stop producer
			return json.Marshal(v)

		case record, ok := <-ch:
			if !ok {
				return json.Marshal(v)
			}

			if len(record.FilePath) > 0 {
				v = append(v, record)
			}
		}
	}
}
