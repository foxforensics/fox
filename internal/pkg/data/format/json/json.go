package json

import (
	"bytes"
	"encoding/json"
)

const (
	prefix = ""
	indent = "  "
)

func Detect(b []byte) bool {
	return json.Valid(b)
}

func Format(b []byte) []byte {
	buf := bytes.NewBuffer(nil)
	err := json.Indent(buf, b, prefix, indent)

	if err != nil {
		return b
	}

	return buf.Bytes()
}
