package hunter

import (
	"os"
	"testing"

	"go.foxforensics.dev/fox/v4/internal/pkg/test"
	"go.foxforensics.dev/fox/v4/internal/pkg/types"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/loader"
)

func TestHunt(t *testing.T) {
	for _, tt := range []struct {
		name  string
		file  string
		count int
	}{
		{
			"empty",
			"misc/empty",
			0,
		}, {
			"evtx",
			"convert/test.evtx.zst",
			919,
		}, {
			"journal",
			"convert/test.journal.zst",
			1922,
		}, {
			"raw",
			"hunt/test.dd.zst",
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

			ldr := loader.New(&loader.Options{
				Limit:    &types.Limits{},
				Filter:   &types.Filters{},
				Parallel: 1,
			})

			file := test.FixtureDeflate(tt.file)

			for range htr.Hunt(ldr.Load([]string{file})) {
				n++
			}

			_ = os.Remove(file)

			if n != tt.count {
				t.Fatal("invalid count:", n)
			}
		})
	}
}
