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

// NoColor toggle
var NoColor = false

// NoSyntax toggle
var NoSyntax = false

var (
	AsGray = gray.SprintFunc()
	AsWarn = warn.SprintFunc()
	AsBold = bold.SprintFunc()
)

var (
	gray = color.New(color.FgHiBlack)
	warn = color.New(color.FgHiRed)
	bold = color.New(color.Bold)
)

var cef = regexp.MustCompile(`[^|]+$`)

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
	if !NoColor && re != nil {
		s = marker.Mark(s, marker.MatchRegexp(re), bold)
	}
	return s
}

func MarkEvent(s string) string {
	if !NoColor {
		s = marker.Mark(s, marker.MatchRegexp(cef), gray)
	}
	return s
}
