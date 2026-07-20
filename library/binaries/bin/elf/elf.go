package elf

import (
	"encoding/json"
	"log/slog"

	"github.com/saferwall/elf"
	"go.foxforensics.eu/fox/v5/library"
)

func Detect(b []byte) bool {
	return library.HasMagic(b, 0, []byte{
		'\x7F', 'E', 'L', 'F',
	})
}

func Convert(b []byte) ([]byte, error) {
	p, err := elf.NewBytes(b)

	if err != nil {
		return b, err
	}

	err = p.Parse()

	if err != nil {
		slog.Warn(err.Error()) // warn only
	}

	return json.Marshal(p.F)
}
