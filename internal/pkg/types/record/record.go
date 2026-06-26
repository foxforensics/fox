package record

import (
	"fmt"
	"strings"

	"go.foxforensics.eu/fox/v4/internal/pkg/types/tables"
	"go.foxforensics.eu/fox/v4/internal/sys/writer"
	"go.foxforensics.eu/hashdump/extract"
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
		s.format(s.LMHash, extract.DefaultLM),
		s.format(s.NTHash, extract.DefaultNT),
	))

	// append historic hashes
	if history {
		for i := range s.NTHashHistory {
			// LM and NT history have always the same length
			sb.WriteString(fmt.Sprintf("\n%s_history%d:%d:%s:%s:::",
				s.SAMAccountName,
				i,
				s.RID,
				s.format(s.LMHashHistory[i], extract.DefaultLM),
				s.format(s.NTHashHistory[i], extract.DefaultNT),
			))
		}
	}

	return sb.String()
}

func (s *Secret) LmOnly() string {
	return s.format(s.LMHash, extract.DefaultLM)
}

func (s *Secret) NtOnly() string {
	return s.format(s.NTHash, extract.DefaultNT)
}

func (s *Secret) format(sum string, def []byte) string {
	if _, pwd := tables.Lookup(sum); len(pwd) > 0 {
		return writer.AsBold(pwd)
	}

	if sum == fmt.Sprintf("%x", def) {
		return writer.AsGray(sum)
	}

	return sum
}
