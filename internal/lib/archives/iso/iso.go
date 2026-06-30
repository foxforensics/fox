package iso

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"strings"

	"github.com/hooklift/iso9660"
	"go.foxforensics.eu/fox/v4/internal/lib"
	"go.foxforensics.eu/fox/v4/internal/sys"
)

func Detect(b []byte) bool {
	for _, o := range []int{
		0x8001, 0x8801, 0x9001,
	} {
		if lib.HasMagic(b, o, []byte{
			'C', 'D', '0', '0', '1',
		}) {
			return true
		}
	}

	return false
}

func Extract(b []byte, root, _ string) (e []lib.Stream) {
	r, err := iso9660.NewReader(bytes.NewReader(b))

	if err != nil {
		slog.Error(err.Error())
		return
	}

	for {
		f, err := r.Next()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			slog.Error(err.Error())
			break
		}

		if f.IsDir() {
			continue
		}

		sr, ok := f.Sys().(io.Reader)

		if !ok {
			slog.Error("invalid type")
			continue
		}

		buf, err := io.ReadAll(sr)

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		e = append(e, lib.Stream{
			Path: sys.JoinPart(root, strings.TrimPrefix(f.Name(), "/")),
			Data: buf,
		})
	}

	return
}
