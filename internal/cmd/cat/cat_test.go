package cat

import (
	"bytes"
	"log/slog"
	"os"
	"testing"

	"go.foxforensics.eu/fox/v5/internal/pkg/tests"
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
				tests.FixtureFile("texts/fox.txt"),
			},
		},
		{
			"Hex",
			"cat.hex.txt",
			[]string{
				"cat",
				"-x",
				tests.FixtureFile("texts/fox.txt"),
			},
		},
		{
			"Auto",
			"cat.auto.txt",
			[]string{
				"cat",
				tests.FixtureFile("binaries/test.mbr"),
			},
		},
		{
			"Json",
			"cat.json",
			[]string{
				"cat",
				tests.FixtureFile("formats/fox.json"),
			},
		},
		{
			"Jsonl",
			"cat.jsonl",
			[]string{
				"cat",
				tests.FixtureFile("formats/fox.jsonl"),
			},
		},
		{
			"Xml",
			"cat.xml",
			[]string{
				"cat",
				tests.FixtureFile("formats/fox.xml"),
			},
		},
		{
			"Find",
			"cat.find.txt",
			[]string{
				"cat",
				"-FO",
				"-C1",
				tests.FixtureFile("texts/fox.txt"),
			},
		},
		{
			"Empty",
			"cat.empty.txt",
			[]string{
				"cat",
				tests.FixtureFile("binaries/test.nil"),
			},
		},
	} {
		for range tests.Cycles {
			t.Run(tt.name, func(t *testing.T) {
				b, err := tests.ExecuteMain(tt.args...)

				if err != nil {
					t.Error(err)
				}

				if !bytes.Equal(b, tests.Sample(tt.sample)) {
					t.Fatal("sample mismatch:", string(b))
				}
			})
		}
	}
}
