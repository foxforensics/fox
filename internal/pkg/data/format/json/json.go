package json

import (
	"bytes"
	"encoding/json"

	"foxhunt.dev/fox/internal/pkg/data/format"
)

func Detect(b []byte) bool {
	return json.Valid(b)
}

func Format(b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	err := json.Indent(buf, b, format.Prefix, format.Indent)

	if err != nil {
		return b, err
	}

	return buf.Bytes(), nil
}
