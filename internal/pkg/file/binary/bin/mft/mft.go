package mft

import (
	"bytes"
	"encoding/json"

	"go.foxforensics.eu/go-mft"

	"go.foxforensics.eu/fox/v4/internal/pkg/file"
)

func Detect(b []byte) bool {
	for _, m := range [][]byte{
		{'F', 'I', 'L', 'E', '*'}, // NT 4.0 & 5.0
		{'F', 'I', 'L', 'E', '0'}, // NT 5.1
	} {
		if file.HasMagic(b, 0, m) {
			return true
		}
	}

	return false
}

func Convert(b []byte) ([]byte, error) {
	dt, err := mft.BuildDirectoryTree(bytes.NewReader(b), "c")

	if err != nil {
		return b, err
	}

	ch := make(chan mft.UsefulMftFields)

	go mft.ParseMftRecords(bytes.NewReader(b), 4096, dt, &ch)

	v := make([]mft.UsefulMftFields, 0, len(b)/1024)

	for record := range ch {
		if len(record.FilePath) > 0 {
			v = append(v, record)
		}
	}

	return json.Marshal(v)
}
