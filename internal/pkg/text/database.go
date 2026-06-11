package text

import (
	"slices"
	"strings"

	"github.com/dlclark/regexp2/v2"
	"go.foxforensics.dev/rhash/database"
)

type Database []database.Entry

func BuildDB(level int) Database {
	var db Database

	if level > 0 {
		db = append(db, []database.Entry{
			{
				Regex: regexp2.MustCompile(`(([a-fA-F0-9]{1,4}:){7}[a-fA-F0-9]{1,4}|([a-fA-F0-9]{1,4}:){1,7}:|([a-fA-F0-9]{1,4}:){1,6}:[a-fA-F0-9]{1,4}|([a-fA-F0-9]{1,4}:){1,5}(:[a-fA-F0-9]{1,4}){1,2}|([a-fA-F0-9]{1,4}:){1,4}(:[a-fA-F0-9]{1,4}){1,3}|([a-fA-F0-9]{1,4}:){1,3}(:[a-fA-F0-9]{1,4}){1,4}|([a-fA-F0-9]{1,4}:){1,2}(:[a-fA-F0-9]{1,4}){1,5}|[a-fA-F0-9]{1,4}:((:[a-fA-F0-9]{1,4}){1,6})|:((:[a-fA-F0-9]{1,4}){1,7}|:)|fe80:(:[a-fA-F0-9]{0,4}){0,4}%[0-9a-zA-Z]+|::(ffff(:0{1,4})?:)?((25[0-5]|(2[0-4]|1?[0-9])?[0-9])\\.){3}(25[0-5]|(2[0-4]|1?[0-9])?[0-9])|([a-fA-F0-9]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1?[0-9])?[0-9])\\.){3}(25[0-5]|(2[0-4]|1?[0-9])?[0-9]))`),
				Names: []string{"IPv6"},
			},
			{
				Regex: regexp2.MustCompile(`\b(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\b`),
				Names: []string{"IPv4"},
			},
			{
				Regex: regexp2.MustCompile(`([a-fA-F0-9]{2}[:-]){5}([a-fA-F0-9]{2})`),
				Names: []string{"MAC"},
			},
			{
				Regex: regexp2.MustCompile("^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$"),
				Names: []string{"DNS"},
			},
			{
				Regex: regexp2.MustCompile("[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,6}"),
				Names: []string{"Mail"},
			},
			{
				Regex: regexp2.MustCompile(`\\\\[a-zA-Z0-9 %._-]+\\[a-zA-Z0-9 $%._-]+`),
				Names: []string{"UNC"},
			},
			{
				Regex: regexp2.MustCompile("[a-zA-Z][a-zA-Z0-9+\\-.]*://([a-zA-Z0-9\\-._~%!$&'()*+,;=]+@)?([a-zA-Z0-9\\-._~%]+|\\[[a-fA-F0-9:.]+]|\\[v[a-fA-F0-9][a-zA-Z0-9\\-._~%!$&'()*+,;=:]+])(:[0-9]+)?(/[a-zA-Z0-9\\-._~%!$&'()*+,;=:@]+)*/?(\\?[a-zA-Z0-9\\-._~%!$&'()*+,;=:@/?]*)?(#[a-zA-Z0-9\\-._~%!$&'()*+,;=:@/?]*)?"),
				Names: []string{"URL"}, // according to RFC 3986
			},
			{
				Regex: regexp2.MustCompile("[a-fA-F0-9]{8}(?:-[a-fA-F0-9]{4}){3}-[a-fA-F0-9]{12}"),
				Names: []string{"UUID"},
			},
			{
				Regex: regexp2.MustCompile(`S-\d-\d+-(\d+-){1,14}\d+`),
				Names: []string{"SID"},
			},
			{
				Regex: regexp2.MustCompile("[0-9]{6}?-[0-9]{6}-[0-9]{6}-[0-9]{6}-[0-9]{6}-[0-9]{6}-[0-9]{6}-[0-9]{6}"),
				Names: []string{"BitLocker"},
			},
			{
				Regex: regexp2.MustCompile(`(HK..|HKEY_[A-Z_]+|\\Registry)\\([a-zA-Z0-9]+\\+)+[a-zA-Z0-9]+`),
				Names: []string{"Registry"},
			},
			{
				Regex: regexp2.MustCompile(`(?:""?[a-zA-Z]:|\\+[^/:*?<>|]+\\+[^/:*?<>|]*)\\+(?:[^/:*?<>|]+\\+)*\w([^/:*?<>|])*`),
				Names: []string{"Windows"}, // https://www.fileside.app/blog/2023-03-17_windows-file-paths/
			},
		}...)
	}

	// https://github.com/EricZimmerman/bstrings/blob/master/bstrings/Program.cs
	if level > 1 {
		db = append(db, []database.Entry{
			{
				Regex: regexp2.MustCompile(`\b[13][a-km-zA-HJ-NP-Z1-9]{25,34}\b`),
				Names: []string{"Bitcoin"},
			},
			{
				Regex: regexp2.MustCompile("Wm[st]{1}[0-9a-zA-Z]{94}"),
				Names: []string{"Aeon"},
			},
			{
				Regex: regexp2.MustCompile("2[0-9AB][0-9a-zA-Z]{93}"),
				Names: []string{"Bytecoin"},
			},
			{
				Regex: regexp2.MustCompile("D[0-9a-zA-Z]{94}"),
				Names: []string{"Dashcoin"},
			},
			{
				Regex: regexp2.MustCompile("(7|X)[a-zA-Z0-9]{33}"),
				Names: []string{"Dashcoin2"},
			},
			{
				Regex: regexp2.MustCompile("6[0-9a-zA-Z]{94}"),
				Names: []string{"Fantomcoin"},
			},
			{
				Regex: regexp2.MustCompile("4[0-9AB][0-9a-zA-Z]{93}|4[0-9AB][0-9a-zA-Z]{104}"),
				Names: []string{"Monero"},
			},
			{
				Regex: regexp2.MustCompile("Sumoo[0-9a-zA-Z]{94}"),
				Names: []string{"Sumokoin"},
			},
			{
				Regex: regexp2.MustCompile(`[ -]*(?:4[ -]*(?:\d[ -]*){11}(?:(?:\d[ -]*){3})?\d|5[ -]*[1-5](?:[ -]*[0-9]){14}|6[ -]*(?:0[ -]*1[ -]*1|5[ -]*\d[ -]*\d)(?:[ -]*[0-9]){12}|3[ -]*[47](?:[ -]*[0-9]){13}|3[ -]*(?:0[ -]*[0-5]|[68][ -]*[0-9])(?:[ -]*[0-9]){11}|(?:2[ -]*1[ -]*3[ -]*1|1[ -]*8[ -]*0[ -]*0|3[ -]*5(?:[ -]*[0-9]){3})(?:[ -]*[0-9]){11})[ -]*`),
				Names: []string{"Credit Card"},
			},
		}...)
	}

	// https://github.com/s0md3v/Bolt/blob/master/db/hashes.json
	if level > 2 {
		db = append(db, database.Entries...)
	}

	return db
}

func (db Database) List() []string {
	var v []string

	for _, e := range db {
		v = append(v, e.Names...)
	}

	slices.SortStableFunc(v, func(a, b string) int {
		return strings.Compare(
			strings.ToLower(a),
			strings.ToLower(b),
		)
	})

	return v
}

func (db Database) Lookup(s string) []string {
	var v []string

	for _, e := range db {
		if ok, _ := e.Regex.MatchString(s); ok {
			v = append(v, e.Names...)
		}
	}

	return v
}
