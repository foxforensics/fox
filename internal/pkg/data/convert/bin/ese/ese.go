package ese

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/Velocidex/ordereddict"
	"www.velocidex.com/golang/go-ese/parser"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

const (
	colAttribId = "ATTc131102"
	colLdapName = "ATTm131532"
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 4, []byte{
		0xEF, 0xCD, 0xAB, 0x89,
	})
}

func Convert(b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	ctx, err := parser.NewESEContext(bytes.NewReader(b))

	if err != nil {
		return b, err
	}

	ctl, err := parser.ReadCatalog(ctx)

	if err != nil {
		return b, err
	}

	rep := strings.NewReplacer(lookup(ctl)...)

	buf.WriteByte('[')

	for i, table := range ctl.Tables.Keys() {
		var json [][]byte

		_ = ctl.DumpTable(table, func(row *ordereddict.Dict) error {
			if b, err := row.MarshalJSON(); err == nil {
				json = append(json, b)
			}
			return nil
		})

		rows := bytes.Join(json, []byte{','})

		buf.WriteString(rep.Replace(fmt.Sprintf(`{"table":"%s","rows":[%s]}`, table, rows)))

		if i < ctl.Tables.Len()-1 {
			buf.WriteByte(',')
		}
	}

	buf.WriteByte(']')

	return buf.Bytes(), nil
}

// https://www.xmco.fr/en/active-directory-en/demystifying-the-ntds-2/
func lookup(ctl *parser.Catalog) []string {
	obj := make(map[int64]string)
	lut := make([]string, 0)

	_ = ctl.DumpTable("MSysObjects", func(row *ordereddict.Dict) error {
		if v, ok := row.GetString("Name"); ok {
			if strings.HasPrefix(v, "ATT") {
				if k, err := strconv.Atoi(v[4:]); err == nil {
					obj[int64(k)] = v
				}
			}
		}
		return nil
	})

	_ = ctl.DumpTable("datatable", func(row *ordereddict.Dict) error {
		if i, ok := row.GetInt64(colAttribId); ok {
			if v, ok := row.GetString(colLdapName); ok {
				if k, ok := obj[i]; ok {
					lut = append(lut, k, v)
				}
			}
		}
		return nil
	})

	return lut
}
