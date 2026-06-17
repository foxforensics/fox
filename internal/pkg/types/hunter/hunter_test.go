package hunter

import (
	"context"
	"testing"

	"go.foxforensics.eu/fox/v4/internal/pkg/test"
	"go.foxforensics.eu/fox/v4/internal/pkg/types"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/loader"
)

func TestHunt(t *testing.T) {
	for _, tt := range []struct {
		name  string
		file  string
		count int
	}{
		{
			"empty",
			"binary/test.nil",
			0,
		}, {
			"evtx",
			"binary/test.evtx",
			3170,
		}, {
			"journal",
			"binary/test.journal",
			1922,
		}, {
			"raw",
			"binary/test.dd",
			919,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var n int

			htr := New(&Options{
				true, // provides also uniqueness
				1,
				0,
			})

			ctx := context.Background()

			ldr := loader.New(&loader.Options{
				Limits:  &types.Limits{},
				Filters: &types.Filters{},
				Threads: 1,
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
