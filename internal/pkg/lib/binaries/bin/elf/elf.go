package elf

import (
	"log/slog"

	"github.com/saferwall/elf"
	"go.foxforensics.eu/fox/v4/internal/pkg/lib"
)

const Magic = "\x7FELF"

func Detect(b []byte) bool {
	return lib.HasMagic(b, 0, []byte(Magic))
}

func Convert(b []byte) ([]byte, error) {
	p, err := elf.NewBytes(b)

	if err != nil {
		return b, err
	}

	err = p.Parse()

	if err != nil {
		slog.Warn(err.Error()) // only warn about missing sections
	}

	raw, err := p.DumpRawJSON()

	if err != nil {
		return b, err
	}

	return []byte(raw), nil
}
