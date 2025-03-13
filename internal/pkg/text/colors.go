package text

import "github.com/fatih/color"

type Colored func(...any) string

var Mark = color.New(color.FgHiWhite).SprintFunc()
var Hide = color.New(color.FgHiBlack).SprintFunc()
var Term = color.New(color.Reset).SprintFunc()
