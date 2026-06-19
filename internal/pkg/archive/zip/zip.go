package zip

import (
	"bytes"
	"io"
	"log/slog"
	"strings"

	"go.foxforensics.eu/fox/v4/internal/pkg"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/go-zip/zip"
)

func Detect(b []byte) bool {
	for _, m := range [][]byte{
		{'P', 'K', 0x03, 0x04}, // default
		{'P', 'K', 0x03, 0x06}, // empty
		{'P', 'K', 0x03, 0x08}, // spanned
	} {
		if pkg.HasMagic(b, 0, m) {
			return true
		}
	}

	return false
}

func Extract(b []byte, root, pass string) (e []pkg.Stream) {
	r, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))

	if err != nil {
		slog.Error(err.Error())
		return
	}

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "/") {
			continue
		}

		if len(pass) > 0 {
			f.SetPassword(pass)
		}

		a, err := f.Open()

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		buf, err := io.ReadAll(a)

		_ = a.Close()

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
