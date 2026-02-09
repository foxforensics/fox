package reg

import (
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/text/encoding/unicode"
	"www.velocidex.com/golang/regparser"
)

func BootKey(r io.ReaderAt) ([]byte, error) {
	var key, b []byte

	reg, err := regparser.NewRegistry(r)

	if err != nil {
		return nil, err
	}

	for _, part := range []string{
		"JD", "Skew1", "GBG", "Data",
	} {
		v := reg.OpenKey(fmt.Sprintf("\\%s\\Control\\Lsa\\%s", controlSet(reg), part))
		a := make([]byte, v.ClassLength())

		_, err = reg.BaseBlock.HiveBin().Reader.ReadAt(a, int64(v.Class()+4096+4))

		if err != nil {
			return nil, err
		}

		b = append(b, a...)
	}

	tmp := string(b)

	if len(b) > 32 {
		dec := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
		tmp, _ = dec.String(string(b))
	}

	sub, err := hex.DecodeString(tmp)

	if err != nil {
		return nil, err
	}

	t := [16]int{8, 5, 4, 2, 11, 9, 13, 3, 0, 6, 1, 12, 14, 10, 15, 7}

	for i := 0; i < len(sub); i++ {
		key = append(key, sub[t[i]])
	}

	return key, nil
}

func controlSet(reg *regparser.Registry) string {
	s := "ControlSet001"

	if k := reg.OpenKey("\\Select"); k != nil {
		for _, v := range k.Values() {
			if v.ValueName() == "Current" {
				s = fmt.Sprintf("ControlSet%03d", v.ValueData().Uint64)
			}
		}
	}

	return s
}
