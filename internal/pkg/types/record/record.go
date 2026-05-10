package record

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.foxforensics.dev/hashdump/extract"

	"go.foxforensics.dev/fox/v4/internal/pkg/text"
)

type Record struct {
	extract.Account
}

func New(base extract.Account) *Record {
	return &Record{base}
}

func (r *Record) ToNTLM(history bool) string {
	var sb strings.Builder

	// append actual hashes
	sb.WriteString(fmt.Sprintf("%s:%d:%s:%s:::",
		r.SAMAccountName,
		r.RID,
		asGray(r.LMHash, extract.DefaultLM),
		asGray(r.NTHash, extract.DefaultNT),
	))

	// append hash histories
	if history {
		for i := range r.NTHashHistory {
			sb.WriteString(fmt.Sprintf("\n%s_history%d:%d:%s:%s:::",
				r.SAMAccountName,
				i,
				r.RID,
				asGray(r.LMHashHistory[i], extract.DefaultLM),
				asGray(r.NTHashHistory[i], extract.DefaultNT),
			))
		}
	}

	return sb.String()
}

func (r *Record) ToJSON() string {
	b, _ := json.MarshalIndent(r, "", "  ")
	return string(b)
}

func (r *Record) ToJSONL() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func asGray(s string, b []byte) string {
	if s == fmt.Sprintf("%x", b) {
		return text.AsGray(s)
	}
	return s
}
