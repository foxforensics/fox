package pf

import (
	"bytes"
	"encoding/json"
	"log/slog"

	"go.foxforensics.eu/fox/v4/internal/pkg/time/entry"
	"go.foxforensics.eu/fox/v4/library"
	"www.velocidex.com/golang/go-prefetch"
)

func Detect(b []byte) bool {
	for _, v := range []struct {
		off int
		buf []byte
	}{
		{off: 4, buf: []byte{'S', 'C', 'C', 'A'}},  // uncompressed
		{off: 0, buf: []byte{'M', 'A', 'M', 0x04}}, // LZX compressed
	} {
		if library.HasMagic(b, v.off, v.buf) {
			return true
		}
	}

	return false
}

func Convert(b []byte) ([]byte, error) {
	pi, err := prefetch.LoadPrefetch(bytes.NewReader(b))

	if err != nil {
		return b, err
	}

	return json.Marshal(pi)
}

func Parse(b []byte) []*entry.Entry {
	var v []*entry.Entry

	pi, err := prefetch.LoadPrefetch(bytes.NewReader(b))

	if err != nil {
		slog.Error(err.Error())
		return v
	}

	for _, t := range pi.LastRunTimes {
		v = append(v, &entry.Entry{
			Name:    pi.Executable,
			Size:    uint64(pi.FileSize),
			Mode:    "-rwxrwxrwx",
			Atime:   t.UTC(),
			Anomaly: pi.RunCount == 0,
		})

		for _, f := range pi.FilesAccessed {
			v = append(v, &entry.Entry{
				Name:  f,
				Mode:  "-rwxrwxrwx",
				Atime: t.UTC(),
			})
		}
	}

	return v
}
