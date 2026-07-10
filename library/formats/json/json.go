package json

import (
	"bytes"
	"encoding/json"

	"go.foxforensics.eu/fox/v5/library/formats"
)

func Detect(b []byte) bool {
	if len(b) > 0 && (b[0] == '{' || b[0] == '[') {
		return json.Valid(b)
	}

	return false
}

func Format(b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	err := json.Indent(buf, b, formats.Prefix, formats.Indent)

	if err != nil {
		return b, err
	}

	return buf.Bytes(), nil
}
