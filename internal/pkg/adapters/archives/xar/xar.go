package xar

import (
	"bytes"
	"io"
	"log/slog"
	"strings"

	"github.com/korylprince/goxar"
	"go.foxforensics.eu/fox/v4/internal/pkg"
	"go.foxforensics.eu/fox/v4/internal/sys"
)

func Detect(b []byte) bool {
	return pkg.HasMagic(b, 0, []byte{
		'x', 'a', 'r', '!',
	})
}

func Extract(b []byte, root, _ string) (e []pkg.Stream) {
	r, err := xar.NewReader(nop(bytes.NewReader(b)), int64(len(b)))

	if err != nil {
		slog.Error(err.Error())
		return
	}

	defer func() {
		_ = r.Close()
	}()

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "/") {
			continue
		}

		if v, err := extractFile(f, root); err == nil {
			e = append(e, v)
		} else {
			slog.Error(err.Error())
		}
	}

	return
}

func extractFile(f *xar.File, root string) (e pkg.Stream, err error) {
	r, err := f.Open()

	if err != nil {
		return e, err
	}

	defer func(r io.Closer) {
		_ = r.Close()
	}(r)

	e.Path = sys.JoinPart(root, f.Name)
	e.Data, err = io.ReadAll(r)

	return e, err
}
