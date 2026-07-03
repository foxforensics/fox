package info

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

func TestInfo(t *testing.T) {
	for _, tt := range []struct {
		name   string
		sample string
		args   []string
	}{
		{
			"Info",
			"info.txt",
			[]string{
				"info",
				test.FixtureFile("binaries"),
			},
		},
		{
			"Json",
			"info.json",
			[]string{
				"info",
				"-j",
				test.FixtureFile("binaries"),
			},
		},
		{
			"Jsonl",
			"info.jsonl",
			[]string{
				"info",
				"-J",
				test.FixtureFile("binaries"),
			},
		},
		{
			"MinMax",
			"info.nx.txt",
			[]string{
				"info",
				"-N6.0",
				"-X7.0",
				test.FixtureFile("binaries"),
			},
		},
		{
			"Block",
			"info.block.txt",
			[]string{
				"info",
				"-B1k",
				test.FixtureFile("binaries/test.rnd"),
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
