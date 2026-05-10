package record

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.foxforensics.dev/hashdump/extract"

	"go.foxforensics.dev/fox/v4/internal/pkg/text"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/rainbow"
)

type Record struct {
	extract.Account
	LMHashCracked string `json:"lm_hash_cracked,omitempty"`
	NTHashCracked string `json:"nt_hash_cracked,omitempty"`
}

func New(account extract.Account) *Record {
	return &Record{
		Account:       account,
		LMHashCracked: rainbow.Lookup(account.LMHash),
		NTHashCracked: rainbow.Lookup(account.NTHash),
	}
}

func (r *Record) ToJSON() string {
	b, _ := json.MarshalIndent(r, "", "  ")
	return string(b)
}

func (r *Record) ToJSONL() string {
	b, _ := json.Marshal(r)
	return string(b)
}

func (r *Record) ToNTLM(history bool) string {
	var sb strings.Builder

	// append actual hashes
	sb.WriteString(fmt.Sprintf("%s:%d:%s:%s:::",
		r.SAMAccountName,
		r.RID,
		format(r.LMHash, extract.DefaultLM),
		format(r.NTHash, extract.DefaultNT),
	))

	// append hash histories
	if history {
		for i := range r.NTHashHistory {
			sb.WriteString(fmt.Sprintf("\n%s_history%d:%d:%s:%s:::",
				r.SAMAccountName,
				i,
				r.RID,
				format(r.LMHashHistory[i], extract.DefaultLM),
				format(r.NTHashHistory[i], extract.DefaultNT),
			))
		}
	}

	return sb.String()
}

func (r *Record) OnlyNTLM() string {
	return fmt.Sprintf("%s:%s",
		format(r.LMHash, extract.DefaultLM),
		format(r.NTHash, extract.DefaultNT),
	)
}

func (r *Record) OnlyLM() string {
	return format(r.LMHash, extract.DefaultLM)
}

func (r *Record) OnlyNT() string {
	return format(r.NTHash, extract.DefaultNT)
}

func format(sum string, def []byte) string {
	if pwd := rainbow.Lookup(sum); len(pwd) > 0 {
		return text.AsWarn(pwd)
	}

	if sum == fmt.Sprintf("%x", def) {
		return text.AsGray(sum)
	}

	return sum
}
