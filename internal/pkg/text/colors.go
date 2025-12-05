package text

import (
	"regexp"

	"github.com/cyucelen/marker"
	"github.com/fatih/color"
)

type Colored func(...any) string

var Mark = white.SprintFunc()
var Hide = black.SprintFunc()
var Term = reset.SprintFunc()

var white = color.New(color.FgHiWhite)
var black = color.New(color.FgHiBlack)
var match = color.New(color.FgHiBlue)
var reset = color.New(color.Reset)

func MarkMatch(s string, re *regexp.Regexp) string {
	return marker.Mark(s, marker.MatchRegexp(re), match)
}
