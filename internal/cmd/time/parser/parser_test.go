package parser

import (
	"context"
	"testing"

	"go.foxforensics.eu/fox/v5/internal/pkg/tests"
)

func TestParse(t *testing.T) {
	for _, tt := range []struct {
		name  string
		file  string
		count int
	}{
		{
			"empty",
			"binaries/test.nil",
			0,
		}, {
			"mft",
			"binaries/test.mft",
			5758,
		}, {
			"lnk",
			"binaries/test.lnk",
			1,
		}, {
			"pf",
			"binaries/test.pf",
			119,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var n int

			for range New(&Options{
				Sort: false,
			}).Parse(context.Background(), tests.Fixture(tt.file)) {
				n++
			}

			if n != tt.count {
				t.Fatalf("invalid count: %d", n)
			}
		})
	}
}
