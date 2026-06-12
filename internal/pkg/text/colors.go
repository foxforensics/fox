package text

import (
	"bytes"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/cyucelen/marker"
	"github.com/dlclark/regexp2/v2"
	"github.com/fatih/color"
)

// Lexer global setting (default)
var Lexer = ""

// Style global setting (default)
var Style = "monokai"

var (
	AsGray = color.New(color.FgHiBlack).SprintfFunc()
	AsWarn = color.New(color.FgHiRed).SprintfFunc()
	AsBold = color.New(color.Bold).SprintfFunc()
)

var fox = color.New(color.Bold).AddRGB(0xff, 0xff, 0xff).AddBgRGB(0x0f, 0x88, 0xcd).Sprint(" FOX ")
var cef = regexp2.MustCompile(`[^|]+$`)

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

func MarkMatch(s string, re *regexp2.Regexp) string {
	if color.NoColor || re == nil {
		return s
	}

	return marker.Mark(s, match(re), color.RGB(0x0f, 0x88, 0xcd))
}

func MarkEvent(s string) string {
	if color.NoColor {
		return s
	}

	return marker.Mark(s, match(cef), color.New(color.FgHiBlack))
}

func match(re *regexp2.Regexp) marker.MatcherFunc {
	return func(s string) marker.Match {
		return marker.Match{
			Template: ReplaceAllString(re, s),
			Patterns: FindAllString(re, s),
		}
	}
}
