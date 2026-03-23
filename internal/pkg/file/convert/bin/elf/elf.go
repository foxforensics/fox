package elf

import (
	"log"
	"strings"

	"github.com/saferwall/elf"

	"github.com/cuhsat/fox/v4/internal/pkg/file"
)

const Magic = "\x7FELF"

func Detect(b []byte) bool {
	return file.HasMagic(b, 0, []byte(Magic))
}

func Convert(b []byte) ([]byte, error) {
	p, err := elf.NewBytes(b)

	if err != nil {
		return b, err
	}

	err = p.Parse()

	if err != nil {
		log.Printf("warning: %s!\n", err)
	}

	raw, err := p.DumpRawJSON()

	if err != nil {
		return []byte(raw), err
	}

	raw = strings.TrimSuffix(raw, "{}")

	return []byte(raw), nil
}
