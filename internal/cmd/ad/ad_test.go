package ad

import (
	"bytes"
	"testing"

	"go.foxforensics.eu/fox/v5/internal/pkg/tests"
)

func TestAd(t *testing.T) {
	for _, tt := range []struct {
		name   string
		sample string
		args   []string
	}{
		{
			"Secrets",
			"ad.secrets.txt",
			[]string{
				"ad",
				tests.FixtureFile("binaries/test.dit"),
				tests.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Secrets Json",
			"ad.secrets.json",
			[]string{
				"ad",
				"-j",
				tests.FixtureFile("binaries/test.dit"),
				tests.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Secrets Jsonl",
			"ad.secrets.jsonl",
			[]string{
				"ad",
				"-l",
				tests.FixtureFile("binaries/test.dit"),
				tests.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Secrets LM",
			"ad.secrets.lm.txt",
			[]string{
				"ad",
				"--lm-only",
				tests.FixtureFile("binaries/test.dit"),
				tests.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Secrets NT",
			"ad.secrets.nt.txt",
			[]string{
				"ad",
				"--nt-only",
				tests.FixtureFile("binaries/test.dit"),
				tests.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Secrets Lookup",
			"ad.lookup.txt",
			[]string{
				"ad",
				"-l",
				tests.FixtureFile("binaries/test.dit"),
				tests.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Computers",
			"ad.computers.txt",
			[]string{
				"ad",
				"-c",
				tests.FixtureFile("binaries/test.dit"),
				tests.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Computers Json",
			"ad.computers.json",
			[]string{
				"ad",
				"-cj",
				tests.FixtureFile("binaries/test.dit"),
				tests.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Computers Jsonl",
			"ad.computers.jsonl",
			[]string{
				"ad",
				"-cJ",
				tests.FixtureFile("binaries/test.dit"),
				tests.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Groups",
			"ad.groups.txt",
			[]string{
				"ad",
				"-g",
				tests.FixtureFile("binaries/test.dit"),
				tests.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Groups Json",
			"ad.groups.json",
			[]string{
				"ad",
				"-gj",
				tests.FixtureFile("binaries/test.dit"),
				tests.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Groups Jsonl",
			"ad.groups.jsonl",
			[]string{
				"ad",
				"-gJ",
				tests.FixtureFile("binaries/test.dit"),
				tests.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Users",
			"ad.users.txt",
			[]string{
				"ad",
				"-u",
				tests.FixtureFile("binaries/test.dit"),
				tests.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Users Json",
			"ad.users.json",
			[]string{
				"ad",
				"-uj",
				tests.FixtureFile("binaries/test.dit"),
				tests.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Users Jsonl",
			"ad.users.jsonl",
			[]string{
				"ad",
				"-uJ",
				tests.FixtureFile("binaries/test.dit"),
				tests.FixtureFile("binaries/test.reg"),
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
