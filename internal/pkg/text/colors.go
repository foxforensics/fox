package text

import (
	"bytes"
	"regexp"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/cyucelen/marker"
	"github.com/fatih/color"
)

// Lexer global setting (default)
var Lexer = ""

// Style global setting (default)
var Style = "monokai"

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

func ColorizeAs(s, hint string) string {
	if color.NoColor {
		return s
	}

	// override hard
	if len(Lexer) > 0 {
		hint = Lexer
	}

	// analyse data
	if len(hint) == 0 {
		if l := lexers.Analyse(s); l != nil {
			hint = l.Config().Name
		}
	}

	if len(hint) == 0 {
		return s
	}

	buf := bytes.NewBuffer(nil)
	err := quick.Highlight(buf, s, hint, "terminal256", Style)

	if err != nil {
		return s
	}

	return buf.String()
}

func MarkMatch(s string, re *regexp.Regexp) string {
	if color.NoColor || re == nil {
		return s
	}

	return marker.Mark(s, marker.MatchRegexp(re), bold)
}

func MarkEvent(s string) string {
	if color.NoColor {
		return s
	}

	return marker.Mark(s, marker.MatchRegexp(cef), gray)
}

func MarkZero(s string) string {
	if color.NoColor {
		return s
	}

	return marker.Mark(s, marker.MatchAll("00"), gray)
}
