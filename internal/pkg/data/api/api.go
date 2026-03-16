package api

import (
	"encoding/json"
)

const (
	Clean       = "clean"
	Unknown     = "unknown"
	Unrated     = "unrated"
	Suspicious  = "suspicious"
	Compromised = "compromised"
)

type Result struct {
	Verdict string            `json:"verdict,omitempty"`
	Details map[string]string `json:"details,omitempty"`
	Stats   struct {
		All int `json:"all,omitempty"`
		Bad int `json:"bad,omitempty"`
	} `json:"stats,omitempty"`
}

func (res *Result) ToJSON() string {
	b, _ := json.MarshalIndent(res, "", "  ")
	return string(b)
}

func (res *Result) ToJSONL() string {
	b, _ := json.Marshal(res)
	return string(b)
}
