package sevenzip

import (
	"bytes"
	"log/slog"

	"github.com/bodgit/sevenzip"
	"go.foxforensics.eu/fox/v5/internal/sys"
	"go.foxforensics.eu/fox/v5/library"
	"go.foxforensics.eu/fox/v5/library/archives"
)

func Detect(b []byte) bool {
	return library.HasMagic(b, 0, []byte{
		0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C,
	})
}

func Extract(b []byte, root, pass string) (e []library.Chunk) {
	r, err := sevenzip.NewReaderWithPassword(bytes.NewReader(b), int64(len(b)), pass)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}

		if v, err := extractFile(f, root, len(b)); err == nil {
			e = append(e, v)
		} else {
			slog.Error(err.Error())
		}
	}

	return
}

// https://github.com/bodgit/sevenzip?tab=readme-ov-file#why-is-my-code-running-so-slow
func extractFile(f *sevenzip.File, root string, size int) (e library.Chunk, err error) {
	r, err := f.Open()

	if err != nil {
		return e, err
	}

	defer func() {
		if err := r.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()

	e.Path = sys.JoinPart(root, f.Name)
	e.Data, err = library.ReadMax(r, size)

	if err != nil {
		err = archives.ErrCorruptPassword
	}

	return e, err
}
