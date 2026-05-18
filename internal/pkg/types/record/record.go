package record

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.foxforensics.dev/hashdump/extract"

	"go.foxforensics.dev/fox/v4/internal/pkg/tables"
	"go.foxforensics.dev/fox/v4/internal/pkg/text"
)

type Record interface {
	String() string
	ToJSON() string
	ToJSONL() string
}

type Computer struct {
	extract.Computer
}

func (c *Computer) String() string {
	return c.Name
}

func (c *Computer) ToJSON() string {
	b, _ := json.MarshalIndent(c, "", "  ")
	return string(b)
}

func (c *Computer) ToJSONL() string {
	b, _ := json.Marshal(c)
	return string(b)
}

type User struct {
	extract.Account
}

func (u *User) String() string {
	return u.Name
}

func (u *User) ToJSON() string {
	b, _ := json.MarshalIndent(u, "", "  ")
	return string(b)
}

func (u *User) ToJSONL() string {
	b, _ := json.Marshal(u)
	return string(b)
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

func (s *Secret) OnlyLM() string {
	return s.format(s.LMHash, extract.DefaultLM)
}

func (s *Secret) OnlyNT() string {
	return s.format(s.NTHash, extract.DefaultNT)
}

func (s *Secret) format(sum string, def []byte) string {
	if pwd := tables.Lookup(sum); len(pwd) > 0 {
		return text.AsBold(pwd)
	}

	if sum == fmt.Sprintf("%x", def) {
		return text.AsGray(sum)
	}

	return sum
}
