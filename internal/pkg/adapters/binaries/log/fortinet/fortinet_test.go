package fortinet

import (
	"strings"
	"testing"

	"go.foxforensics.eu/fox/v4/internal/test"
)

const src = "binaries/test.tlog"

func BenchmarkDetect(b *testing.B) {
	buf := test.Fixture(src)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkConvert(b *testing.B) {
	buf := test.Fixture(src)

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
			if !Detect(test.Fixture(tt.path)) {
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
			buf, err := Convert(test.Fixture(tt.path))

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
