package elf

import (
	"log/slog"
	"strings"

	"github.com/saferwall/elf"
	"go.foxforensics.eu/fox/v4/internal/pkg"
)

const Magic = "\x7FELF"

func Detect(b []byte) bool {
	return pkg.HasMagic(b, 0, []byte(Magic))
}

func Convert(b []byte) ([]byte, error) {
	p, err := elf.NewBytes(b)

	if err != nil {
		return b, err
	}

	err = p.Parse()

	if err != nil {
		slog.Warn(err.Error())
	}

	raw, err := p.DumpRawJSON()

	if err != nil {
		return []byte(raw), err
	}

	raw = strings.TrimSuffix(raw, "{}")

	return []byte(raw), nil
}
