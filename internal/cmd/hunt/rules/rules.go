package rules

import (
	_ "embed"
	"slices"
	"strings"

	"github.com/bradleyjkemp/sigma-go"
)

//go:embed rules.yml
var Critical []byte

func IsSupported(r *sigma.Rule) bool {
	return slices.Contains([]string{
		"fox",
		"linux",
		"windows",
	}, strings.ToLower(r.Logsource.Product))
}
