package cab

import (
	"bytes"
	"log/slog"
	"strings"

	"github.com/secDre4mer/go-cab"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/fox/v4/library"
)

func Detect(b []byte) bool {
	return library.HasMagic(b, 0, []byte{
		0x4D, 0x53, 0x43, 0x46,
	})
}

func Extract(b []byte, root, _ string) (e []library.Chunk) {
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

		buf, err := library.ReadMax(rc, len(b))

		if err := rc.Close(); err != nil {
			slog.Error(err.Error())
		}

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		e = append(e, library.Chunk{
			Path: sys.JoinPart(root, f.Name),
			Data: buf,
		})
	}

	return
}
