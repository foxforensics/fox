package cat

import (
	"bytes"
	"log/slog"
	"os"
	"testing"

	"go.foxforensics.eu/fox/v4/internal/test"
)

func TestMain(m *testing.M) {
	if err := os.Chdir("../../../"); err != nil {
		slog.Error(err.Error())
	} else {
		os.Exit(m.Run())
	}
}

func TestCat(t *testing.T) {
	for _, tt := range []struct {
		name   string
		sample string
		args   []string
	}{
		{
			"Text",
			"cat.text.txt",
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
			"Auto",
			"cat.auto.txt",
			[]string{
				"cat",
				test.FixtureFile("binaries/test.mbr"),
			},
		},
		{
			"Json",
			"cat.json",
			[]string{
				"cat",
				test.FixtureFile("formats/fox.json"),
			},
		},
		{
			"Jsonl",
			"cat.jsonl",
			[]string{
				"cat",
				test.FixtureFile("formats/fox.jsonl"),
			},
		},
		{
			"Xml",
			"cat.xml",
			[]string{
				"cat",
				test.FixtureFile("formats/fox.xml"),
			},
		},
		{
			"Find",
			"cat.find.txt",
			[]string{
				"cat",
				"-FO",
				"-C1",
				test.FixtureFile("texts/fox.txt"),
			},
		},
		{
			"Empty",
			"cat.empty.txt",
			[]string{
				"cat",
				test.FixtureFile("binaries/test.nil"),
			},
		},
	} {
		for range test.Cycles {
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
}
