package hunter

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/data/reader/ewf"
	"github.com/cuhsat/fox/v4/internal/pkg/data/reader/vhdx"
	"github.com/cuhsat/fox/v4/internal/pkg/data/reader/vmdk"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
	"github.com/cuhsat/fox/v4/internal/pkg/types/loader"
	"github.com/cuhsat/fox/v4/internal/pkg/types/register"
)

func TestMain(m *testing.M) {
	register.Reader("ewf", ewf.Detect, ewf.Reader)
	register.Reader("vhdx", vhdx.Detect, vhdx.Reader)
	register.Reader("vmdk", vmdk.Detect, vmdk.Reader)

	os.Exit(m.Run())
}

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
			"convert/test.evtx",
			919,
		}, {
			"journal",
			"convert/test.journal",
			1922,
		}, {
			"raw",
			"hunt/nist.dd",
			17336,
		}, {
			"ewf",
			"hunt/nist.E01",
			818,
		}, {
			"vhdx",
			"hunt/nist.vhdx",
			17336,
		}, {
			"vmdk",
			"hunt/nist.vmdk",
			17336,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			htr := New(&Options{
				true, // provides also uniqueness
				1,
				0,
			})

			events := consume(htr, fixture(tt.file))

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

func fixture(file string) string {
	const dir = "../../../../testdata"

	_, c, _, ok := runtime.Caller(0)

	if !ok {
		log.Fatalln("runtime error")
	}

	return filepath.Join(filepath.Dir(c), dir, file)
}
