package lnk

import (
	"bytes"
	"encoding/json"

	"go.foxforensics.eu/go-lnk"

	"go.foxforensics.eu/fox/v4/internal/pkg/file"
)

func Detect(b []byte) bool {
	return file.HasMagic(b, 0, []byte{
		0x4C, 0, 0, 0,
	})
}

func Convert(b []byte) ([]byte, error) {
	lf, err := lnk.Read(bytes.NewReader(b), uint64(len(b)))

	if err != nil {
		return b, err
	}

	return json.Marshal(lf)
}
