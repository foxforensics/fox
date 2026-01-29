package ese

import (
	"bytes"

	"github.com/Velocidex/ordereddict"
	"www.velocidex.com/golang/go-ese/parser"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
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

	for _, t := range ctl.Tables.Keys() {
		if err = ctl.DumpTable(t, func(row *ordereddict.Dict) error {
			row.Set("table", t)

			b, err := row.MarshalJSON()

			if err != nil {
				return err
			}

			buf.Write(b)
			buf.WriteByte('\n')

			return nil
		}); err != nil {
			return buf.Bytes(), err
		}
	}

	return buf.Bytes(), nil
}
