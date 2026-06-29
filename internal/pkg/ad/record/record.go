package record

import (
	"fmt"
	"strings"

	"go.foxforensics.eu/fox/v4/internal/pkg/ad/tables"
	"go.foxforensics.eu/fox/v4/internal/sys/writer"
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
	return c.Name
}

type Group struct {
	extract.Group
}

func (g *Group) String() string {
	return g.Name
}

type User struct {
	extract.Account
}

func (u *User) String() string {
	return u.Name
}

type Secret struct {
	extract.Account
}

func (s *Secret) ToNTLM(history bool) string {
	var sb strings.Builder

	// append actual hashes
	sb.WriteString(fmt.Sprintf("%s:%d:%s:%s:::",
		s.SAMAccountName,
		s.RID,
		s.format(s.LMHash, defaultLM),
		s.format(s.NTHash, defaultNT),
	))

	// append historic hashes
	if history {
		for i := range s.NTHashHistory {
			var lm = defaultLM
			var nt = defaultNT

			if i < len(s.LMHashHistory) {
				lm = s.LMHashHistory[i]
			}

			if i < len(s.NTHashHistory) {
				nt = s.NTHashHistory[i]
			}

			sb.WriteString(fmt.Sprintf("\n%s_history%d:%d:%s:%s:::",
				s.SAMAccountName,
				i,
				s.RID,
				s.format(lm, defaultLM),
				s.format(nt, defaultNT),
			))
		}
	}

	return sb.String()
}

func (s *Secret) LmOnly() string {
	return s.format(s.LMHash, defaultLM)
}

func (s *Secret) NtOnly() string {
	return s.format(s.NTHash, defaultNT)
}

func (s *Secret) format(sum, def string) string {
	if _, pwd := tables.Lookup(sum); len(pwd) > 0 {
		return writer.AsBold(pwd)
	}

	if sum == def {
		return writer.AsGray(sum)
	}

	return sum
}
