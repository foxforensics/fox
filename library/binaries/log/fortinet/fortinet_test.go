package fortinet

import (
	"strings"
	"testing"

	"go.foxforensics.eu/fox/v5/internal/pkg/tests"
)

const src = "binaries/test.tlog"

func BenchmarkDetect(b *testing.B) {
	buf := tests.Fixture(src)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkConvert(b *testing.B) {
	buf := tests.Fixture(src)

	for b.Loop() {
		_, _ = Convert(buf)
	}
}

func TestDetect(t *testing.T) {
	for _, tt := range []struct {
		name string
		path string
	}{
		{
			"elog",
			"binaries/test.elog",
		},
		{
			"tlog",
			"binaries/test.tlog",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if !Detect(tests.Fixture(tt.path)) {
				t.Fatal("not detected")
			}
		})
	}
}

func TestConvert(t *testing.T) {
	for _, tt := range []struct {
		name string
		path string
	}{
		{
			"elog",
			"binaries/test.elog",
		},
		{
			"tlog",
			"binaries/test.tlog",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			buf, err := Convert(tests.Fixture(tt.path))

			if err != nil {
				t.Error(err)
			}

			lines := strings.Split(string(buf), "\n")

			if len(lines) == 0 {
				t.Fatal("invalid length")
			}
		})
	}
}
