package color

import (
	"bytes"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
)

const style = "monokai"

func Detect(b []byte) bool {
	return lexers.Analyse(string(b)) != nil
}

func Format(b []byte) []byte {
	buf := bytes.NewBuffer(nil)
	err := quick.Highlight(buf, string(b), "", "terminal256", style)

	if err != nil {
		return b
	}

	return buf.Bytes()
}
