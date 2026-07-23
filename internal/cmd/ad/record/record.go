package record

import (
	"fmt"
	"strings"

	"go.foxforensics.eu/fox/v5/internal/cmd/ad/tables"
	"go.foxforensics.eu/fox/v5/internal/pkg"
	"go.foxforensics.eu/fox/v5/internal/pkg/writer"
	"go.foxforensics.eu/hashdump/extract"
)

var (
	defaultLM = fmt.Sprintf("%x", extract.DefaultLM)
	defaultNT = fmt.Sprintf("%x", extract.DefaultNT)
)

type Record interface {
	String() string
}

type Computer struct {
	extract.Computer
}

func (c *Computer) String() string {
	return pkg.Sanitize(c.Name)
}

type Group struct {
	extract.Group
}

func (g *Group) String() string {
	return pkg.Sanitize(g.Name)
}

type User struct {
	extract.Account
}

func (u *User) String() string {
	return pkg.Sanitize(u.Name)
}

type Secret struct {
	extract.Account
}

func (s *Secret) ToNTLM(history bool) string {
	var sb strings.Builder

	// append actual hashes
	fmt.Fprintf(&sb, "%s:%d:%s:%s:::",
		pkg.Sanitize(s.SAMAccountName),
		s.RID,
		s.format(s.LMHash, defaultLM, true),
		s.format(s.NTHash, defaultNT, false),
	)

	// append historic hashes
	if history {
		for i := range max(len(s.LMHashHistory), len(s.NTHashHistory)) {
			var lm = defaultLM
			var nt = defaultNT

			if i < len(s.LMHashHistory) {
				lm = s.LMHashHistory[i]
			}

			if i < len(s.NTHashHistory) {
				nt = s.NTHashHistory[i]
			}

			fmt.Fprintf(&sb, "\n%s_history%d:%d:%s:%s:::",
				pkg.Sanitize(s.SAMAccountName),
				i,
				s.RID,
				s.format(lm, defaultLM, true),
				s.format(nt, defaultNT, false),
			)
		}
	}

	return sb.String()
}

func (s *Secret) LmOnly() string {
	return s.format(s.LMHash, defaultLM, true)
}

func (s *Secret) NtOnly() string {
	return s.format(s.NTHash, defaultNT, false)
}

func (s *Secret) format(sum, def string, upper bool) string {
	if _, pwd := tables.Lookup(sum); len(pwd) > 0 {
		if upper {
			pwd = strings.ToUpper(pwd)
		}
		return writer.AsBold(pwd)
	}

	if sum == def {
		return writer.AsGray(sum)
	}

	return sum
}
