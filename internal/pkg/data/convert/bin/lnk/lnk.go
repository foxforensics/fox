package lnk

import (
	"bytes"
	"encoding/json"

	"github.com/cuhsat/golnk"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte{
		0x4c, 0, 0, 0,
	})
}

func Convert(b []byte) ([]byte, error) {
	lf, err := lnk.Read(bytes.NewReader(b), uint64(len(b)))

	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(lf, "", "  ")
}
