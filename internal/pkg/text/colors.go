package text

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/cyucelen/marker"
	"github.com/fatih/color"
)

const Lexer = "json"
const Style = "monokai"

var Hide = black.SprintFunc()
var Warn = alert.SprintFunc()

var black = color.New(color.FgHiBlack)
var match = color.New(color.BgBlue)
var alert = color.New(color.FgHiRed)

func MarkMatch(s string, re *regexp.Regexp) string {
	if re == nil {
		return s // no regex, no match
	}

	return marker.Mark(s, marker.MatchRegexp(re), match)
}

func ColorizeAs(s, l string) string {
	var sb strings.Builder

	if color.NoColor {
		return s
	}

	err := quick.Highlight(&sb, s, l, "terminal256", Style)

	if err != nil {
		return s
	}

	return sb.String()
}

func Colorize(b []byte) []byte {
	if color.NoColor {
		return b
	}

	s := string(b)

	lexer := Lexer // use fallback

	if l := lexers.Analyse(s); l != nil {
		lexer = l.Config().Name
	}

	buf := bytes.NewBuffer(nil)

	err := quick.Highlight(buf, s, lexer, "terminal256", Style)

	if err != nil {
		return b
	}

	return buf.Bytes()
}
