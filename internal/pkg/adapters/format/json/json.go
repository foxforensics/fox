package json

import (
	"bytes"
	"encoding/json"

	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/format"
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
