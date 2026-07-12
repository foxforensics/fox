package time

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
				"-s",
				tests.FixtureFile("binaries/test.mft"),
			},
		},
		{
			"Json",
			"time.json",
			[]string{
				"time",
				"-j",
				tests.FixtureFile("binaries/test.lnk"),
			},
		},
		{
			"Jsonl",
			"time.jsonl",
			[]string{
				"time",
				"-J",
				tests.FixtureFile("binaries/test.lnk"),
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
