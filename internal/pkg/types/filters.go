package types

import (
	"regexp"

	"github.com/cuhsat/fox/v4/internal/pkg/types/smap"
)

type Filters struct {
	Regex  *regexp.Regexp // regex pattern
	Before uint           // lines before
	After  uint           // lines after
}

func (f *Filters) FilterSMap(s smap.SMap) smap.SMap {
	if f.Regex == nil {
		return s // not filtered
	}

	v := s.Grep(f.Regex)

	if f.Before+f.After == 0 {
		return v // without context
	}

	r := make(smap.SMap, 0, len(v))

	for grp, str := range v {
		for _, b := range (s)[max((str.Nr-1)-f.Before, 0) : str.Nr-1] {
			b.Grp = uint(grp + 1)
			r = append(r, b)
		}

		str.Grp = uint(grp + 1)
		r = append(r, str)

		for _, a := range (s)[str.Nr:min(str.Nr+f.After, uint(len(s)))] {
			a.Grp = uint(grp + 1)
			r = append(r, a)
		}
	}

	return r // with context
}
