package ese

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/Velocidex/ordereddict"
	"go.foxforensics.eu/fox/v5/library"
	"go.foxforensics.eu/go-ese/parser"
)

// row attributes
const (
	attribId = "ATTc131102"
	ldapName = "ATTm131532"
)

var wrapper = `{"table":%s,"rows":[%s]}`

func Detect(b []byte) bool {
	return library.HasMagic(b, 4, []byte{
		0xEF, 0xCD, 0xAB, 0x89,
	})
}

func Convert(b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	ctx, err := parser.NewESEContext(bytes.NewReader(b), int64(len(b)))

	if err != nil {
		return b, err
	}

	ctl, err := parser.ReadCatalog(ctx)

	if err != nil {
		return b, err
	}

	rep := strings.NewReplacer(translate(ctl)...)

	buf.WriteByte('[')

	for i, table := range ctl.Tables.Keys() {
		var v [][]byte

		if err = ctl.DumpTable(table, func(row *ordereddict.Dict) error {
			if b, err := marshal(row, rep); err == nil {
				v = append(v, b)
			}
			return nil
		}); err != nil {
			slog.Warn(err.Error())
		}

		rows := bytes.Join(v, []byte(","))

		name, err := json.Marshal(table)

		if err != nil {
			name = []byte(err.Error())
		}

		fmt.Fprintf(buf, wrapper, name, rows)

		if i < ctl.Tables.Len()-1 {
			buf.WriteByte(',')
		}
	}

	buf.WriteByte(']')

	return buf.Bytes(), nil
}

// translate attributes to their LDAP name
// source: https://www.xmco.fr/en/active-directory-en/demystifying-the-ntds-2/
func translate(ctl *parser.Catalog) []string {
	cache := make(map[int64]string)
	names := make([]string, 0)

	// build name attribute cache
	if err := ctl.DumpTable("MSysObjects", func(row *ordereddict.Dict) error {
		if v, ok := row.GetString("Name"); ok && len(v) > 4 {
			if strings.HasPrefix(v, "ATT") {
				if k, err := strconv.Atoi(v[4:]); err == nil {
					cache[int64(k)] = v
				}
			}
		}
		return nil
	}); err != nil {
		slog.Warn(err.Error())
	}

	// generate translations
	if err := ctl.DumpTable("datatable", func(row *ordereddict.Dict) error {
		if i, ok := row.GetInt64(attribId); ok {
			if v, ok := row.GetString(ldapName); ok {
				if k, ok := cache[i]; ok {
					names = append(names, k, v)
				}
			}
		}
		return nil
	}); err != nil {
		slog.Warn(err.Error())
	}

	return names
}

func marshal(od *ordereddict.Dict, rep *strings.Replacer) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte('{')

	for _, item := range od.Items() {
		key, err := json.Marshal(rep.Replace(item.Key))

		if err != nil {
			continue
		}

		// Check for back references and skip them
		sub, ok := item.Value.(*ordereddict.Dict)

		if ok && sub == od {
			continue
		}

		buf.Write(key)
		buf.WriteByte(':')

		val, err := json.Marshal(item.Value)

		if err == nil {
			buf.Write(val)
			buf.WriteByte(',')
		} else {
			buf.WriteString("null,")
		}
	}

	if buf.Len() > 1 {
		buf.Truncate(buf.Len() - 1)
	}

	buf.WriteByte('}')

	return buf.Bytes(), nil
}
