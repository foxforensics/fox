package ad

import (
	"bytes"
	"testing"

	"go.foxforensics.eu/fox/v4/internal/test"
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
				test.FixtureFile("binaries/test.dit"),
				test.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Secrets Json",
			"ad.secrets.json",
			[]string{
				"ad",
				"-j",
				test.FixtureFile("binaries/test.dit"),
				test.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Secrets Jsonl",
			"ad.secrets.jsonl",
			[]string{
				"ad",
				"-J",
				test.FixtureFile("binaries/test.dit"),
				test.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Secrets LM",
			"ad.secrets.lm.txt",
			[]string{
				"ad",
				"--lm-only",
				test.FixtureFile("binaries/test.dit"),
				test.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Secrets NT",
			"ad.secrets.nt.txt",
			[]string{
				"ad",
				"--nt-only",
				test.FixtureFile("binaries/test.dit"),
				test.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Secrets Lookup",
			"ad.lookup.txt",
			[]string{
				"ad",
				"-l",
				test.FixtureFile("binaries/test.dit"),
				test.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Computers",
			"ad.computers.txt",
			[]string{
				"ad",
				"-c",
				test.FixtureFile("binaries/test.dit"),
				test.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Computers Json",
			"ad.computers.json",
			[]string{
				"ad",
				"-cj",
				test.FixtureFile("binaries/test.dit"),
				test.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Computers Jsonl",
			"ad.computers.jsonl",
			[]string{
				"ad",
				"-cJ",
				test.FixtureFile("binaries/test.dit"),
				test.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Groups",
			"ad.groups.txt",
			[]string{
				"ad",
				"-g",
				test.FixtureFile("binaries/test.dit"),
				test.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Groups Json",
			"ad.groups.json",
			[]string{
				"ad",
				"-gj",
				test.FixtureFile("binaries/test.dit"),
				test.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Groups Jsonl",
			"ad.groups.jsonl",
			[]string{
				"ad",
				"-gJ",
				test.FixtureFile("binaries/test.dit"),
				test.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Users",
			"ad.users.txt",
			[]string{
				"ad",
				"-u",
				test.FixtureFile("binaries/test.dit"),
				test.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Users Json",
			"ad.users.json",
			[]string{
				"ad",
				"-uj",
				test.FixtureFile("binaries/test.dit"),
				test.FixtureFile("binaries/test.reg"),
			},
		},
		{
			"Users Jsonl",
			"ad.users.jsonl",
			[]string{
				"ad",
				"-uJ",
				test.FixtureFile("binaries/test.dit"),
				test.FixtureFile("binaries/test.reg"),
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
