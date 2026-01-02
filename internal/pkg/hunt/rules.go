package hunt

import (
	"strings"

	"github.com/bradleyjkemp/sigma-go"
)

var Default = []byte(`
title: Fox Hunt Windows
id: f0badbad-ff01-40af-9b83-f1a74aef8174
description: Detects critical Windows system events
author: Christian Uhsat [mail@foxhunt.wtf]
status: stable
logsource:
  product: windows
detection:
  selection:
    EventID:
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
level: critical
`)

func IsCompatible(r *sigma.Rule) bool {
	return strings.ToLower(r.Logsource.Product) == "windows"
}
