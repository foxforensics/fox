package hunt

import (
	"bytes"
	"log/slog"
	"os"
	"testing"

	"go.foxforensics.eu/fox/v5/internal/pkg/tests"
)

type Request struct {
	body        []byte
	authHeader  string
	contentType string
}

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
				tests.FixtureFile("binaries/test.dd"),
			},
		},
		{
			"Json",
			"hunt.json",
			[]string{
				"hunt",
				"-saj",
				tests.FixtureFile("binaries/test.dd"),
			},
		},
		{
			"Jsonl",
			"hunt.jsonl",
			[]string{
				"hunt",
				"-saJ",
				tests.FixtureFile("binaries/test.dd"),
			},
		},
		{
			"Triage",
			"hunt.triage.txt",
			[]string{
				"hunt",
				"-t",
				tests.FixtureFile("binaries/test.evtx"),
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

func TestStream(t *testing.T) {
	for _, tt := range []struct {
		name   string
		sample string
		args   []string
	}{
		{
			"Triage",
			"hunt.triage.txt",
			[]string{
				"hunt",
				"-t",
				tests.FixtureFile("binaries/test.evtx"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tests.ExecuteMain(tt.args...)

			if err != nil {
				t.Error(err)
			}
		})
	}
}
