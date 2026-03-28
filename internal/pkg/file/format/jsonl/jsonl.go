package jsonl

import (
	"bufio"
	"bytes"
	"encoding/json"

	"go.foxforensics.dev/fox/v4/internal/pkg/file/format"
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

func Format(b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	r := bufio.NewReader(bytes.NewReader(b))

	for {
		line, _, err := r.ReadLine()

		if err != nil {
			break
		}

		err = json.Indent(buf, line, format.Prefix, format.Indent)

		if err != nil {
			buf.Write(line)
		}

		buf.WriteByte('\n')
	}

	return buf.Bytes(), nil
}
