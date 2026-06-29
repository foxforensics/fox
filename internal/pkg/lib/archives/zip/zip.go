package zip

import (
	"bytes"
	"io"
	"log/slog"
	"strings"

	"go.foxforensics.eu/fox/v4/internal/pkg/lib"
	"go.foxforensics.eu/fox/v4/internal/pkg/lib/archives"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/go-zip/zip"
)

func Detect(b []byte) bool {
	for _, m := range [][]byte{
		{'P', 'K', 0x03, 0x04}, // default
		{'P', 'K', 0x03, 0x06}, // empty
	} {
		if lib.HasMagic(b, 0, m) {
			return true
		}
	}

	return false
}

func Extract(b []byte, root, pass string) (e []lib.Stream) {
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

		rc, err := saveOpen(f)

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

func saveOpen(f *zip.File) (rc io.ReadCloser, err error) {
	defer func() {
		if r := recover(); r != nil {
			rc, err = nil, archives.ErrCorruptPassword
		}
	}()

	rc, err = f.Open()

	return
}
