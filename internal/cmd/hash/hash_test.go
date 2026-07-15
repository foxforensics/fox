package hash

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

func TestHash(t *testing.T) {
	for _, tt := range []struct {
		name   string
		sample string
		args   []string
	}{
		{
			"Hash",
			"hash.txt",
			[]string{
				"hash",
				tests.FixtureFile("binaries"),
			},
		},
		{
			"Json",
			"hash.json",
			[]string{
				"hash",
				"-j",
				tests.FixtureFile("binaries"),
			},
		},
		{
			"Jsonl",
			"hash.jsonl",
			[]string{
				"hash",
				"-l",
				tests.FixtureFile("binaries"),
			},
		},
		{
			"All",
			"hash.all.txt",
			[]string{
				"hash",
				"-a",
				tests.FixtureFile("binaries/test.rnd"),
			},
		},
		{
			"Exe",
			"hash.exe.txt",
			[]string{
				"hash",
				"-Himpfuzzy,imphasho,imphashs,sdhash,tlsh",
				tests.FixtureFile("binaries/test??.exe"),
			},
		},
		{
			"Img",
			"hash.img.txt",
			[]string{
				"hash",
				"-Haverage,blockmean,difference,marrhildreth,median,pdq,phash,rash,whash",
				tests.FixtureFile("images/test.jpg"),
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
