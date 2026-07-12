package str

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

func TestStr(t *testing.T) {
	for _, tt := range []struct {
		name   string
		sample string
		args   []string
	}{
		{
			"Str",
			"str.txt",
			[]string{
				"str",
				tests.FixtureFile("texts/test.txt"),
			},
		},
		{
			"Ascii",
			"str.ascii.txt",
			[]string{
				"str",
				"-a",
				tests.FixtureFile("texts/test.txt"),
			},
		},
		{
			"Sort",
			"str.sort.txt",
			[]string{
				"str",
				"-s",
				tests.FixtureFile("texts/test.txt"),
			},
		},
		{
			"Trim",
			"str.trim.txt",
			[]string{
				"str",
				"-at",
				tests.FixtureFile("texts/test.txt"),
			},
		},
		{
			"MinMax",
			"str.nx.txt",
			[]string{
				"str",
				"-N6",
				"-X8",
				tests.FixtureFile("texts/test.txt"),
			},
		},
		{
			"What1",
			"str.what1.txt",
			[]string{
				"str",
				"-w",
				tests.FixtureFile("texts/test.txt"),
			},
		},
		{
			"What2",
			"str.what2.txt",
			[]string{
				"str",
				"-ww",
				tests.FixtureFile("texts/test.txt"),
			},
		},
		{
			"What3",
			"str.what3.txt",
			[]string{
				"str",
				"-www",
				tests.FixtureFile("texts/test.txt"),
			},
		},
		{
			"Class",
			"str.class.txt",
			[]string{
				"str",
				"-Cipv4,ipv6",
				tests.FixtureFile("texts/test.txt"),
			},
		},
	} {
		for range tests.Cycles {
			t.Run(tt.name, func(t *testing.T) {
				b, err := tests.FixtureMain(tt.args...)

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
