package info

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
				tests.FixtureFile("binaries"),
			},
		},
		{
			"Json",
			"info.json",
			[]string{
				"info",
				"-j",
				tests.FixtureFile("binaries"),
			},
		},
		{
			"Jsonl",
			"info.jsonl",
			[]string{
				"info",
				"-l",
				tests.FixtureFile("binaries"),
			},
		},
		{
			"MinMax",
			"info.nx.txt",
			[]string{
				"info",
				"-N6.0",
				"-X7.0",
				tests.FixtureFile("binaries"),
			},
		},
		{
			"Block",
			"info.block.txt",
			[]string{
				"info",
				"-B1k",
				tests.FixtureFile("binaries/test.rnd"),
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
