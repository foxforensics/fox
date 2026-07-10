package time

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

func TestTime(t *testing.T) {
	for _, tt := range []struct {
		name   string
		sample string
		args   []string
	}{
		{
			"Time",
			"time.txt",
			[]string{
				"time",
				test.FixtureFile("binaries/test.mft"),
			},
		},
		{
			"Json",
			"time.json",
			[]string{
				"time",
				"-j",
				test.FixtureFile("binaries/test.lnk"),
			},
		},
		{
			"Jsonl",
			"time.jsonl",
			[]string{
				"time",
				"-J",
				test.FixtureFile("binaries/test.lnk"),
			},
		},
		{
			"Bodyfile",
			"time.bodyfile.csv",
			[]string{
				"time",
				"-b",
				test.FixtureFile("binaries/test.pf"),
			},
		},
		{
			"Timesketch",
			"time.timesketch.jsonl",
			[]string{
				"time",
				"-t",
				test.FixtureFile("binaries/test.pf"),
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
