package hunter

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heap"
)

func TestHunt(t *testing.T) {
	for _, tt := range []struct {
		name string
		file string
		cnt  int
	}{
		{
			"EventLogs",
			"convert/test.evtx",
			919,
		}, {
			"Journals",
			"convert/test.journal",
			1922,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			l := New(newOpts())

			events := consume(l, fixture(tt.file))

			if len(events) != tt.cnt {
				t.Fatal("invalid count")
			}
		})
	}
}

func newOpts() *Options {
	return &Options{
		true,
		3,
		1,
		0,
	}
}

func newCtx() *heap.Context {
	return &heap.Context{
		Name:   "test",
		Type:   types.Regular,
		Limit:  &types.Limits{},
		Filter: &types.Filters{},
	}
}

func consume(htr *Hunter, b []byte) (out []*event.Event) {
	h := heap.New(newCtx(), b)

	ch := make(chan *heap.Heap, 1)
	ch <- h

	close(ch)

	for evt := range htr.Hunt(ch) {
		out = append(out, evt)
	}

	h.Discard()

	return out
}

func fixture(file string) []byte {
	const dir = "../../../../testdata"

	_, c, _, ok := runtime.Caller(0)

	if !ok {
		log.Fatalln("runtime error")
	}

	b, err := os.ReadFile(filepath.Join(filepath.Dir(c), dir, file))

	if err != nil {
		log.Fatalln(err)
	}

	return b
}
