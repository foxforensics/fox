package api

import (
	"encoding/hex"
	"encoding/json"

	"github.com/xxtea/xxtea-go/xxtea"
)

const (
	Clean      = "clean"
	Unknown    = "unknown"
	Unrated    = "unrated"
	Suspicious = "suspicious"
)

type Report struct {
	Verdict string            `json:"verdict,omitempty"`
	Details map[string]string `json:"details,omitempty"`
	Stats   struct {
		All int `json:"all,omitempty"`
		Bad int `json:"bad,omitempty"`
	} `json:"stats,omitempty"`
}

func (rep *Report) ToJSON() string {
	b, _ := json.MarshalIndent(rep, "", "  ")
	return string(b)
}

func (rep *Report) ToJSONL() string {
	b, _ := json.Marshal(rep)
	return string(b)
}

func Decrypt(s, k string) string {
	v, _ := hex.DecodeString(s)
	return string(xxtea.Decrypt(v, []byte(k)))
}
