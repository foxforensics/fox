// Package xml source: https://github.com/golang/go/issues/58994
package xml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"

	"github.com/cuhsat/fox/v4/internal/pkg/file"
	"github.com/cuhsat/fox/v4/internal/pkg/file/format"
)

func Detect(b []byte) bool {
	return file.HasMagic(b, 0, []byte{
		'<', '?', 'x', 'm', 'l', ' ',
	})
}

func Format(b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	err := indent(buf, bytes.NewReader(b), format.Prefix, format.Indent)

	if err != nil {
		return b, err
	}

	return bytes.Replace(buf.Bytes(), []byte("?>"), []byte("?>\n"), 1), nil
}

func indent(w io.Writer, r io.Reader, prefix, indent string) error {
	decode := xml.NewDecoder(r)
	encode := xml.NewEncoder(w)
	encode.Indent(prefix, indent)

	for {
		token, err := decode.Token()

		if errors.Is(err, io.EOF) {
			return encode.Flush()
		} else if err != nil {
			return err
		}

		if data, ok := token.(xml.CharData); ok {
			token = xml.CharData(bytes.TrimSpace(data))
		}

		if err := encode.EncodeToken(token); err != nil {
			return err
		}
	}
}
