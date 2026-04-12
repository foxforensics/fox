package record

import (
	"encoding/json"
	"fmt"

	"go.foxforensics.dev/hashdump/pkg/hashdump"

	"go.foxforensics.dev/fox/v4/internal/pkg/text"
)

type Record struct {
	hashdump.Record
}

func New(base hashdump.Record) *Record {
	return &Record{base}
}

func (r *Record) String() string {
	lm := r.LM
	nt := r.NT

	if lm == fmt.Sprintf("%x", hashdump.EmptyLM) {
		lm = text.AsGray(lm)
	}

	if nt == fmt.Sprintf("%x", hashdump.EmptyNT) {
		nt = text.AsGray(nt)
	}

	return fmt.Sprintf("%s:%d:%s:%s:::", r.Username, r.RID, lm, nt)
}

func (r *Record) ToJSON() string {
	b, _ := json.MarshalIndent(r, "", "  ")
	return string(b)
}

func (r *Record) ToJSONL() string {
	b, _ := json.Marshal(r)
	return string(b)
}
