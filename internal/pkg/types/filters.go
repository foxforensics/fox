package types

import (
	"github.com/dlclark/regexp2/v2"

	"go.foxforensics.dev/fox/v4/internal/pkg/types/smap"
)

type Filters struct {
	Regex  *regexp2.Regexp // regex pattern
	Before uint            // lines before
	After  uint            // lines after
}

func (f *Filters) Filter(s smap.SMap) smap.SMap {
	if f.Regex == nil {
		return s // not filtered
	}

	v := s.Grep(f.Regex)

	if f.Before+f.After == 0 {
		return v // without context
	}

	r := make(smap.SMap, 0, len(v))

	for grp, str := range v {
		for _, b := range (s)[max(int((str.Line-1)-f.Before), 0) : str.Line-1] {
			b.Group = uint(grp + 1)
			r = append(r, b)
		}

		str.Group = uint(grp + 1)
		r = append(r, str)

		for _, a := range (s)[str.Line:min(int(str.Line+f.After), len(s))] {
			a.Group = uint(grp + 1)
			r = append(r, a)
		}
	}

	return r // with context
}
