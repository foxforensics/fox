package pe

import (
	"encoding/json"

	"github.com/saferwall/pe"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

const Magic = "MZ"

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte(Magic))
}

func Format(b []byte) ([]byte, error) {
	p, err := pe.NewBytes(b, new(pe.Options))

	if err != nil {
		return nil, err
	}

	err = p.Parse()

	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(p, "", "  ")
}
