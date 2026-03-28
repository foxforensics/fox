package fortinet

import (
	"bytes"

	"github.com/f0x4n6/go-fortilog/pkg/fortilog"

	"github.com/cuhsat/fox/v4/internal/pkg/file"
)

func Detect(b []byte) bool {
	for _, m := range [][]byte{
		{0xEC, 0xCE}, // llog v5 old format
		{0xEC, 0xCF}, // llog v5 old format
		{0xEC, 0xDF}, // llog v5 new format
		{0xAA, 0x01}, // llog v5 tlc block
	} {
		if file.HasMagic(b, 0, m) {
			return true
		}
	}

	return false
}

func Convert(b []byte) ([]byte, error) {
	var buf bytes.Buffer

	err := fortilog.DecodeLLogV5(b, &buf)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
