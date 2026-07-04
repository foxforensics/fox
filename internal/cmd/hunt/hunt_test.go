package hunt

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

func TestHunt(t *testing.T) {
	for _, tt := range []struct {
		name   string
		sample string
		args   []string
	}{
		{
			"Hunt",
			"hunt.txt",
			[]string{
				"hunt",
				"-sa",
				test.FixtureFile("binaries/test.dd"),
			},
		},
		{
			"Json",
			"hunt.json",
			[]string{
				"hunt",
				"-saj",
				test.FixtureFile("binaries/test.dd"),
			},
		},
		{
			"Jsonl",
			"hunt.jsonl",
			[]string{
				"hunt",
				"-saJ",
				test.FixtureFile("binaries/test.dd"),
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
