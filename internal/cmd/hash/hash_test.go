package hash

import (
	"bytes"
	"log/slog"
	"os"
	"testing"

	"go.foxforensics.eu/fox/v5/internal/test"
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
				test.FixtureFile("binaries"),
			},
		},
		{
			"Json",
			"hash.json",
			[]string{
				"hash",
				"-j",
				test.FixtureFile("binaries"),
			},
		},
		{
			"Jsonl",
			"hash.jsonl",
			[]string{
				"hash",
				"-J",
				test.FixtureFile("binaries"),
			},
		},
		{
			"All",
			"hash.all.txt",
			[]string{
				"hash",
				"-a",
				test.FixtureFile("binaries/test.rnd"),
			},
		},
		{
			"Exe",
			"hash.exe.txt",
			[]string{
				"hash",
				"-Himpfuzzy,imphasho,imphashs,sdhash,tlsh",
				test.FixtureFile("binaries/test??.exe"),
			},
		},
		{
			"Img",
			"hash.img.txt",
			[]string{
				"hash",
				"-Haverage,blockmean,difference,marrhildreth,median,pdq,phash,rash,whash",
				test.FixtureFile("images/test.jpg"),
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
					t.Fatal("sample mismatch:", string(b))
				}
			})
		}
	}
}
