package carver

import (
	"context"
	"testing"

	"go.foxforensics.eu/fox/v5/internal/test"
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
			"binaries/test.nil",
			false,
			0,
		}, {
			"strings",
			"texts/test.txt",
			false,
			14,
		}, {
			"nasty",
			"texts/nasty.txt",
			false,
			582,
		}, {
			"exe",
			"binaries/fox.exe",
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
