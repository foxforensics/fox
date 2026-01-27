package jsonl

import (
	"bufio"
	"bytes"
	"encoding/json"

	"github.com/cuhsat/fox/v4/internal/pkg/data/format"
)

func Detect(b []byte) bool {
	r := bufio.NewReader(bytes.NewReader(b))

	line, _, err := r.ReadLine()

	if err != nil {
		return false
	}

	// test only the first line for performance reasons
	return json.Valid(line)
}

func Format(b []byte) []byte {
	buf := bytes.NewBuffer(nil)

	r := bufio.NewReader(bytes.NewReader(b))

	for {
		b, _, err := r.ReadLine()

		if err != nil {
			break
		}

		err = json.Indent(buf, b, format.Prefix, format.Indent)

		if err != nil {
			break
		}

		buf.WriteByte('\n')
	}

	return buf.Bytes()
}
