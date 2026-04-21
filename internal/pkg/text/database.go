package text

import (
	"regexp"
	"slices"
	"strings"

	"go.foxforensics.dev/rhash/database"
)

type Database []database.Entry

func BuildDB(level int) Database {
	var db Database

	if level > 0 {
		db = append(db, []database.Entry{
			{
				regexp.MustCompile(`(([a-fA-F0-9]{1,4}:){7}[a-fA-F0-9]{1,4}|([a-fA-F0-9]{1,4}:){1,7}:|([a-fA-F0-9]{1,4}:){1,6}:[a-fA-F0-9]{1,4}|([a-fA-F0-9]{1,4}:){1,5}(:[a-fA-F0-9]{1,4}){1,2}|([a-fA-F0-9]{1,4}:){1,4}(:[a-fA-F0-9]{1,4}){1,3}|([a-fA-F0-9]{1,4}:){1,3}(:[a-fA-F0-9]{1,4}){1,4}|([a-fA-F0-9]{1,4}:){1,2}(:[a-fA-F0-9]{1,4}){1,5}|[a-fA-F0-9]{1,4}:((:[a-fA-F0-9]{1,4}){1,6})|:((:[a-fA-F0-9]{1,4}){1,7}|:)|fe80:(:[a-fA-F0-9]{0,4}){0,4}%[0-9a-zA-Z]+|::(ffff(:0{1,4})?:)?((25[0-5]|(2[0-4]|1?[0-9])?[0-9])\\.){3}(25[0-5]|(2[0-4]|1?[0-9])?[0-9])|([a-fA-F0-9]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1?[0-9])?[0-9])\\.){3}(25[0-5]|(2[0-4]|1?[0-9])?[0-9]))`),
				[]string{"IPv6"},
			},
			{
				regexp.MustCompile(`(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])`),
				[]string{"IPv4"},
			},
			{
				regexp.MustCompile(`([a-fA-F0-9]{2}[:-]){5}([a-fA-F0-9]{2})`),
				[]string{"MAC"},
			},
			{
				regexp.MustCompile("[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,6}"),
				[]string{"Mail"},
			},
			{
				regexp.MustCompile(`\\\\[a-zA-Z0-9 %._-]+\\[a-zA-Z0-9 $%._-]+`),
				[]string{"UNC"},
			},
			{
				regexp.MustCompile("[a-zA-Z][a-zA-Z0-9+\\-.]*://([a-zA-Z0-9\\-._~%!$&'()*+,;=]+@)?([a-zA-Z0-9\\-._~%]+|\\[[a-fA-F0-9:.]+]|\\[v[a-fA-F0-9][a-zA-Z0-9\\-._~%!$&'()*+,;=:]+])(:[0-9]+)?(/[a-zA-Z0-9\\-._~%!$&'()*+,;=:@]+)*/?(\\?[a-zA-Z0-9\\-._~%!$&'()*+,;=:@/?]*)?(#[a-zA-Z0-9\\-._~%!$&'()*+,;=:@/?]*)?"),
				[]string{"URL"}, // according to RFC 3986
			},
			{
				regexp.MustCompile("[a-fA-F0-9]{8}(?:-[a-fA-F0-9]{4}){3}-[a-fA-F0-9]{12}"),
				[]string{"UUID"},
			},
			{
				regexp.MustCompile(`S-\d-\d+-(\d+-){1,14}\d+`),
				[]string{"SID"},
			},
			{
				regexp.MustCompile("[0-9]{6}?-[0-9]{6}-[0-9]{6}-[0-9]{6}-[0-9]{6}-[0-9]{6}-[0-9]{6}-[0-9]{6}"),
				[]string{"BitLocker"},
			},
			{
				regexp.MustCompile(`(HK..|HKEY_[A-Z_]+|\\Registry)\\([a-zA-Z0-9]+\\+)+[a-zA-Z0-9]+`),
				[]string{"Registry"},
			},
			{
				regexp.MustCompile(`(?:""?[a-zA-Z]:|\\+[^/:*?<>|]+\\+[^/:*?<>|]*)\\+(?:[^/:*?<>|]+\\+)*\w([^/:*?<>|])*`),
				[]string{"Windows"}, // https://www.fileside.app/blog/2023-03-17_windows-file-paths/
			},
		}...)
	}

	// https://github.com/EricZimmerman/bstrings/blob/master/bstrings/Program.cs
	if level > 1 {
		db = append(db, []database.Entry{
			{
				regexp.MustCompile(`[13][a-km-zA-HJ-NP-Z1-9]{25,34}`),
				[]string{"Bitcoin"},
			},
			{
				regexp.MustCompile("Wm[st][0-9a-zA-Z]{94}"),
				[]string{"Aeon"},
			},
			{
				regexp.MustCompile("2[0-9AB][0-9a-zA-Z]{93}"),
				[]string{"Bytecoin"},
			},
			{
				regexp.MustCompile("D[0-9a-zA-Z]{94}"),
				[]string{"Dashcoin"},
			},
			{
				regexp.MustCompile("([7X])[a-zA-Z0-9]{33}"),
				[]string{"Dashcoin2"},
			},
			{
				regexp.MustCompile("6[0-9a-zA-Z]{94}"),
				[]string{"Fantomcoin"},
			},
			{
				regexp.MustCompile("4[0-9AB][0-9a-zA-Z]{93}|4[0-9AB][0-9a-zA-Z]{104}"),
				[]string{"Monero"},
			},
			{
				regexp.MustCompile("Sumoo[0-9a-zA-Z]{94}"),
				[]string{"Sumokoin"},
			},
			{
				regexp.MustCompile(`4[0-9]{12}(?:[0-9]{3})?|[25][1-7][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\d{3})\d{11}`),
				[]string{"Credit Card"},
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
		if e.Regex.MatchString(s) {
			v = append(v, e.Names...)
		}
	}

	return v
}
