package text

import "github.com/dlclark/regexp2/v2"

// ReplaceAllString is a regexp compatibility function.
func ReplaceAllString(re *regexp2.Regexp, s string) string {
	v, _ := re.Replace(s, "%s", 0, -1)
	return v
}

// FindAllString is a regexp compatibility function.
func FindAllString(re *regexp2.Regexp, s string) []string {
	var v []string
	m, _ := re.FindStringMatch(s)
	for m != nil {
		v = append(v, m.String())
		m, _ = re.FindNextMatch(m)
	}
	return v
}
