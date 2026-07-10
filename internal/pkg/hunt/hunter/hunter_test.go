package hunter

import (
	"context"
	"testing"

	"go.foxforensics.eu/fox/v5/internal/pkg"
	"go.foxforensics.eu/fox/v5/internal/sys/loader"
	"go.foxforensics.eu/fox/v5/internal/test"
)

func TestHunt(t *testing.T) {
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
			"evtx",
			"binaries/test.evtx",
			3170,
		}, {
			"journal",
			"binaries/test.journal",
			1922,
		}, {
			"raw",
			"binaries/test.dd",
			919,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var n int

			htr := New(&Options{
				true, // provides also uniqueness
				1,
			})

			ctx := context.Background()

			ldr := loader.New(&loader.Options{
				Query: pkg.Query{},
			})

			file := test.FixtureFile(tt.file)

			for range htr.Hunt(ctx, ldr.Load(ctx, []string{file})) {
				n++
			}

			if n != tt.count {
				t.Fatal("invalid count:", n)
			}
		})
	}
}
