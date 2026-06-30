package cab

import (
	"bytes"
	"io"
	"log/slog"
	"strings"

	"github.com/secDre4mer/go-cab"
	"go.foxforensics.eu/fox/v4/internal/lib"
	"go.foxforensics.eu/fox/v4/internal/sys"
)

func Detect(b []byte) bool {
	return lib.HasMagic(b, 0, []byte{
		0x4D, 0x53, 0x43, 0x46,
	})
}

func Extract(b []byte, root, _ string) (e []lib.Stream) {
	r, err := cab.Open(bytes.NewReader(b), int64(len(b)))

	if err != nil {
		slog.Error(err.Error())
		return
	}

	for _, f := range r.Files {
		if strings.HasSuffix(f.Name, "/") {
			continue
		}

		rc, err := f.Open()

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		buf, err := io.ReadAll(rc)

		if err := rc.Close(); err != nil {
			slog.Error(err.Error())
		}

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		e = append(e, lib.Stream{
			Path: sys.JoinPart(root, f.Name),
			Data: buf,
		})
	}

	return
}
