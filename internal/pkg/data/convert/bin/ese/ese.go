package ese

import (
	"bytes"
	"encoding/json"

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
	buf := bytes.NewBuffer(b)

	ctx, err := parser.NewESEContext(bytes.NewReader(b))

	if err != nil {
		return nil, err
	}

	ctl, err := parser.ReadCatalog(ctx)

	if err != nil {
		return nil, err
	}

	for _, t := range ctl.Tables.Keys() {
		if err = ctl.DumpTable(t, func(row *ordereddict.Dict) error {
			b, err := row.MarshalJSON()

			if err != nil {
				return err
			}

			// sanity check
			if json.Valid(b) {
				buf.Write(b)
				buf.WriteByte('\n')

				println(string(b))
			} else {
				return nil
			}

			return nil
		}); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
