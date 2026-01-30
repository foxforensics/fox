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

const Lexer = "plain"
const Style = "monokai"

var Hide = black.SprintFunc()
var Warn = alert.SprintFunc()

var black = color.New(color.FgHiBlack)
var match = color.New(color.BgBlue)
var alert = color.New(color.FgHiRed)

var mapping = map[string]string{
	"elf":     "json",
	"exe":     "json",
	"dll":     "json",
	"sys":     "json",
	"ese":     "json",
	"evtx":    "json",
	"journal": "json",
	"lnk":     "json",
	"pf":      "json",
	"json":    "json",
	"jsonl":   "json",
	"txt":     "plain",
}

func MarkMatch(s string, re *regexp.Regexp) string {
	if re == nil {
		return s // no regex, no match
	}

	return marker.Mark(s, marker.MatchRegexp(re), match)
}

func ColorizeStringAs(s, lexer string) string {
	var sb strings.Builder

	if color.NoColor {
		return s
	}

	err := quick.Highlight(&sb, s, lexer, "terminal256", Style)

	if err != nil {
		return s
	}

	return sb.String()
}

func ColorizeAs(b []byte, lexer string) []byte {
	buf := bytes.NewBuffer(nil)

	if color.NoColor {
		return b
	}

	err := quick.Highlight(buf, string(b), lexer, "terminal256", Style)

	if err != nil {
		return b
	}

	return buf.Bytes()
}

func Colorize(b []byte, hint string) []byte {
	var lexer string

	if v, ok := mapping[hint]; ok {
		lexer = v // use type mapping
	} else if l := lexers.Analyse(string(b)); l != nil {
		lexer = l.Config().Name
	} else {
		lexer = Lexer // use fallback
	}

	return ColorizeAs(b, lexer)
}
