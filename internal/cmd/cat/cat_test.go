package cat

import (
	"bytes"
	"testing"

	"go.foxforensics.eu/fox/v4/internal/test"
)

func TestCat(t *testing.T) {
	for _, tt := range []struct {
		name   string
		sample string
		args   []string
	}{
		{
			"Text",
			"cat.txt",
			[]string{
				"cat",
				"-t",
				test.FixtureFile("texts/fox.txt"),
			},
		},
		{
			"Hex",
			"cat.hex.txt",
			[]string{
				"cat",
				"-x",
				test.FixtureFile("texts/fox.txt"),
			},
		},
		{
			"Find",
			"cat.find.txt",
			[]string{
				"cat",
				"-FO -C1",
				test.FixtureFile("texts/fox.txt"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			b, err := test.FixtureMain(tt.args...)

			if err != nil {
				t.Error(err)
			}

			if !bytes.Equal(b, test.Sample(tt.sample)) {
				t.Fatal("sample mismatch")
			}
		})
	}
}
