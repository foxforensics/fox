package text

import (
	"bytes"
	"regexp"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/cyucelen/marker"
	"github.com/fatih/color"
)

// Lexer (default)
const Lexer = "text"

// Style (default)
const Style = "monokai"

var NoSyntax = false

var (
	AsBlue = color.RGB(0xff, 0xff, 0xff).AddBgRGB(0x0f, 0x88, 0xcd).SprintFunc()
	AsGray = color.New(color.FgHiBlack).SprintFunc()
	AsWarn = color.New(color.FgHiRed).SprintFunc()
	AsBold = color.New(color.Bold).SprintFunc()
)

func ColorizeStringAs(s, lexer string) string {
	return string(ColorizeAs([]byte(s), lexer, Style))
}

func ColorizeAs(b []byte, lexer string, style string) []byte {
	if NoSyntax {
		return b
	}

	if style == "" {
		style = Style
	}

	buf := bytes.NewBuffer(nil)

	err := quick.Highlight(buf, string(b), lexer, "terminal256", style)

	if err != nil {
		return b
	}

	return buf.Bytes()
}

func Colorize(b []byte, hint string, style string) []byte {
	if NoSyntax {
		return b
	}

	var lexer string

	if len(hint) > 0 {
		lexer = hint // use type hinting
	} else if l := lexers.Analyse(string(b)); l != nil {
		lexer = l.Config().Name
	} else {
		lexer = Lexer // use fallback
	}

	if lexer == "text" {
		return b
	}

	return ColorizeAs(b, lexer, style)
}

func MarkMatch(s string, re *regexp.Regexp) string {
	if re == nil {
		return s // no regex, no match
	}

	return marker.Mark(s, marker.MatchRegexp(re), color.New(color.Bold))
}
