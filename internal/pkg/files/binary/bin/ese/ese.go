package ese

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/Velocidex/ordereddict"
	"go.foxforensics.eu/fox/v4/internal/pkg"
	"go.foxforensics.eu/go-ese/parser"
)

// row attributes
const (
	attribId = "ATTc131102"
	ldapName = "ATTm131532"
)

// table wrapper
var wrapper = "{\"table\":\"%s\",\"rows\":[\n%s\n]}"

func Detect(b []byte) bool {
	return pkg.HasMagic(b, 4, []byte{
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
		var json [][]byte

		_ = ctl.DumpTable(table, func(row *ordereddict.Dict) error {
			if b, err := row.MarshalJSON(); err == nil {
				json = append(json, b)
			}
			return nil
		})

		rows := bytes.Join(json, []byte(",\n"))

		buf.WriteString(rep.Replace(fmt.Sprintf(wrapper, table, rows)))

		if i < ctl.Tables.Len()-1 {
			buf.WriteString(",\n")
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
	_ = ctl.DumpTable("MSysObjects", func(row *ordereddict.Dict) error {
		if v, ok := row.GetString("Name"); ok {
			if strings.HasPrefix(v, "ATT") {
				if k, err := strconv.Atoi(v[4:]); err == nil {
					cache[int64(k)] = v
				}
			}
		}
		return nil
	})

	// generate translations
	_ = ctl.DumpTable("datatable", func(row *ordereddict.Dict) error {
		if i, ok := row.GetInt64(attribId); ok {
			if v, ok := row.GetString(ldapName); ok {
				if k, ok := cache[i]; ok {
					names = append(names, k, v)
				}
			}
		}
		return nil
	})

	return names
}
