package cab

import (
	"bytes"
	"io"
	"log/slog"
	"strings"

	"github.com/secDre4mer/go-cab"
	"go.foxforensics.eu/fox/v4/internal/pkg"
	"go.foxforensics.eu/fox/v4/internal/sys"
)

func Detect(b []byte) bool {
	return pkg.HasMagic(b, 0, []byte{
		0x4D, 0x53, 0x43, 0x46,
	})
}

func Extract(b []byte, root, _ string) (e []pkg.Stream) {
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

		_ = rc.Close()

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		e = append(e, pkg.Stream{
			Path: sys.JoinPart(root, f.Name),
			Data: buf,
		})
	}

	return
}
