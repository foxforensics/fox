package json

import (
	"bytes"
	"encoding/json"

	"github.com/cuhsat/fox/v4/internal/pkg/data/format"
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
