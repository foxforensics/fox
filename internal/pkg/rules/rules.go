package rules

import (
	"slices"
	"strings"

	"github.com/bradleyjkemp/sigma-go"
)

var Default = []byte(`
title: Fox Standard Rule
logsource:
  product: fox
detection: 
  selection:
    - PRIORITY:
      - 0
      - 1
      - 2
      - 3
    - EventID:
      - 1015
      - 1102
      - 1116
      - 1117
      - 1119
      - 4624
      - 4625
      - 4647
      - 4648
      - 4663
      - 4672
      - 4697
      - 4719
      - 4728
      - 4732
      - 4735
      - 4740
      - 4756
      - 4771
      - 4776
      - 4820
      - 4821
      - 4822
      - 4823
      - 4824
      - 4964
  condition: selection
`)

var supported = []string{
	"fox",
	"linux",
	"windows",
}

func IsSupported(r *sigma.Rule) bool {
	return slices.Contains(supported, strings.ToLower(r.Logsource.Product))
}
