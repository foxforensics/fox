package hunt

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
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
				"-sal",
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
	var h http.Header

	buf := bytes.NewBuffer(nil)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// cache bodies
		_, _ = buf.ReadFrom(r.Body)
		buf.WriteByte('\n')

		// cache last header
		h = r.Header

		w.WriteHeader(http.StatusOK)
	}))

	defer srv.Close()

	for _, tt := range []struct {
		name   string
		mime   string
		token  string
		sample string
		args   []string
	}{
		{
			"CEF",
			"text/plain",
			"",
			"hunt.cef.txt",
			[]string{
				"hunt",
				"-sU",
				srv.URL,
				tests.FixtureFile("binaries/test.evtx"),
			},
		},
		{
			"ECS",
			"application/json",
			"",
			"hunt.ecs.jsonl",
			[]string{
				"hunt",
				"-sE",
				srv.URL,
				tests.FixtureFile("binaries/test.evtx"),
			},
		},
		{
			"HEC",
			"application/json",
			"Splunk Test",
			"hunt.hec.jsonl",
			[]string{
				"hunt",
				"-sH",
				srv.URL,
				"-A",
				"Test",
				tests.FixtureFile("binaries/test.evtx"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tests.ExecuteMain(tt.args...)

			if err != nil {
				t.Error(err)
			}

			if !bytes.Equal(buf.Bytes(), tests.Sample(tt.sample)) {
				t.Fatal("sample mismatch")
			}

			if h.Get("Authorization") != tt.token {
				t.Fatal("token mismatch")
			}

			if h.Get("Content-Type") != tt.mime {
				t.Fatal("type mismatch")
			}

			buf.Reset()
		})
	}
}
