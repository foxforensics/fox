package mft

import (
	"bytes"
	"encoding/json"

	"github.com/AlecRandazzo/MFT-Parser"

	"go.foxforensics.dev/fox/v4/internal/pkg/file"
)

func Detect(b []byte) bool {
	return file.HasMagic(b, 0, []byte{
		'F', 'I', 'L', 'E',
	})
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

	return json.MarshalIndent(v, "", "  ")
}
