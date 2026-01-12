package text

import (
	"regexp"

	"github.com/cyucelen/marker"
	"github.com/fatih/color"
)

type Colored func(...any) string

var Mark = white.SprintFunc()
var Hide = black.SprintFunc()
var Warn = alert.SprintFunc()
var Term = reset.SprintFunc()

var white = color.New(color.FgHiWhite)
var black = color.New(color.FgHiBlack)
var match = color.New(color.FgHiBlue)
var alert = color.New(color.FgHiRed)
var reset = color.New(color.Reset)

func MarkMatch(s string, re *regexp.Regexp) string {
	if re == nil {
		return s // no regex, no match
	}

	return marker.Mark(s, marker.MatchRegexp(re), match)
}

func MarkMatchFunc(re *regexp.Regexp) Colored {
	return func(a ...any) string {
		return MarkMatch(a[0].(string), re)
	}
}
