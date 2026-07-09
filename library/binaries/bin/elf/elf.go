package elf

import (
	"log/slog"
	"strings"

	"github.com/saferwall/elf"
	"go.foxforensics.eu/fox/v4/library"
)

const Magic = "\x7FELF"

func Detect(b []byte) bool {
	return library.HasMagic(b, 0, []byte(Magic))
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

	// the marshaled JSON is invalid, so we must wrap it
	raw = `{"binary":` + strings.Replace(raw, "}{", `},"symbols":{`, 1) + "}"

	return []byte(raw), nil
}
