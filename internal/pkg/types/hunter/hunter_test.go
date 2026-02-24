package hunter

import (
	"os"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/test"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
	"github.com/cuhsat/fox/v4/internal/pkg/types/loader"
)

func TestHunt(t *testing.T) {
	for _, tt := range []struct {
		name string
		file string
		cnt  int
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
			htr := New(&Options{
				true, // provides also uniqueness
				1,
				0,
			})

			file := test.FixtureDeflate(tt.file)

			events := consume(htr, file)

			defer func() {
				_ = os.Remove(file)
			}()

			if len(events) != tt.cnt {
				t.Fatal("invalid count:", len(events))
			}
		})
	}
}

func consume(htr *Hunter, path string) (out []*event.Event) {
	ldr := loader.New(&loader.Options{
		Limit:    &types.Limits{},
		Filter:   &types.Filters{},
		Parallel: 1,
	})

	for evt := range htr.Hunt(ldr.Load([]string{path})) {
		out = append(out, evt)
	}

	return out
}
