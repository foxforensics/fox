package str

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
				test.FixtureFile("texts/test.txt"),
			},
		},
		{
			"Ascii",
			"str.ascii.txt",
			[]string{
				"str",
				"-a",
				test.FixtureFile("texts/test.txt"),
			},
		},
		{
			"Sort",
			"str.sort.txt",
			[]string{
				"str",
				"-s",
				test.FixtureFile("texts/test.txt"),
			},
		},
		{
			"Trim",
			"str.trim.txt",
			[]string{
				"str",
				"-at",
				test.FixtureFile("texts/test.txt"),
			},
		},
		{
			"MinMax",
			"str.nx.txt",
			[]string{
				"str",
				"-N6",
				"-X8",
				test.FixtureFile("texts/test.txt"),
			},
		},
		{
			"What1",
			"str.what1.txt",
			[]string{
				"str",
				"-w",
				test.FixtureFile("texts/test.txt"),
			},
		},
		{
			"What2",
			"str.what2.txt",
			[]string{
				"str",
				"-ww",
				test.FixtureFile("texts/test.txt"),
			},
		},
		{
			"What3",
			"str.what3.txt",
			[]string{
				"str",
				"-www",
				test.FixtureFile("texts/test.txt"),
			},
		},
		{
			"Class",
			"str.class.txt",
			[]string{
				"str",
				"-Cipv4,ipv6",
				test.FixtureFile("texts/test.txt"),
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
