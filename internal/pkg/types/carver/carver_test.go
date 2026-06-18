package carver

import (
	"context"
	"testing"

	"go.foxforensics.eu/fox/v4/internal/pkg/test"
)

func TestCarve(t *testing.T) {
	for _, tt := range []struct {
		name  string
		file  string
		ascii bool
		count int
	}{
		{
			"empty",
			"binary/test.nil",
			false,
			0,
		}, {
			"strings",
			"string/test.txt",
			false,
			14,
		}, {
			"nasty",
			"string/nasty.txt",
			false,
			582,
		}, {
			"exe",
			"binary/fox.exe",
			true,
			15446,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var n int

			crv := New(&Options{
				Min:   4,
				Max:   256,
				Ascii: tt.ascii,
			})

			for range crv.Carve(context.Background(), test.Fixture(tt.file)) {
				n++
			}

			if n != tt.count {
				t.Fatalf("invalid count: %d", n)
			}
		})
	}
}
