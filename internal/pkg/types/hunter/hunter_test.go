package hunter

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/xz"
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
			"convert/test.evtx.xz",
			919,
		}, {
			"Journals",
			"convert/test.journal.xz",
			1922,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			htr := New(newOpts())

			events := consume(htr, fixture(tt.file))

			if len(events) != tt.cnt {
				t.Fatal("invalid count")
			}
		})
	}
}

func newOpts() *Options {
	return &Options{
		true,
		1,
		0,
	}
}

func consume(htr *Hunter, b []byte) (out []*event.Event) {
	h := heap.New("test", b, new(types.Limits))

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

	buf, err := os.ReadFile(filepath.Join(filepath.Dir(c), dir, file))

	if err != nil {
		log.Fatalln(err)
	}

	if !xz.Detect(buf) {
		return buf
	}

	buf, err = xz.Deflate(buf)

	if err != nil {
		log.Fatalln(err)
	}

	return buf
}
