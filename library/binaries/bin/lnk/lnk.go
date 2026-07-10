package lnk

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"

	"go.foxforensics.eu/fox/v5/internal/pkg/time/entry"
	"go.foxforensics.eu/fox/v5/library"
	"go.foxforensics.eu/go-lnk"
)

func Detect(b []byte) bool {
	return library.HasMagic(b, 0, []byte{
		0x4C, 0, 0, 0,
	})
}

func Convert(b []byte) ([]byte, error) {
	lf, err := lnk.Read(bytes.NewReader(b), uint64(len(b)))

	if err != nil {
		return b, err
	}

	return json.Marshal(lf)
}

func Parse(b []byte) []entry.Entry {
	var v []entry.Entry

	lf, err := lnk.Read(bytes.NewReader(b), uint64(len(b)))

	if err != nil {
		slog.Error(err.Error())
		return v
	}

	return append(v, entry.Entry{
		Name:  lf.LinkInfo.LocalBasePath + lf.LinkInfo.CommonPathSuffix,
		Mode:  buildMode(&lf),
		Size:  uint64(lf.Header.TargetFileSize),
		Mtime: lf.Header.WriteTime.UTC(),
		Atime: lf.Header.AccessTime.UTC(),
		Btime: lf.Header.CreationTime.UTC(),
	})
}

func buildMode(lf *lnk.File) string {
	var sb strings.Builder

	if v, ok := lf.Header.FileAttributes["FILE_ATTRIBUTE_DIRECTORY"]; ok && v {
		sb.WriteByte('d')
	} else if v, ok := lf.Header.FileAttributes["REPARSE_POINT"]; ok && v {
		sb.WriteByte('l')
	} else {
		sb.WriteByte('-')
	}

	if v, ok := lf.Header.FileAttributes["FILE_ATTRIBUTE_READONLY"]; ok && v {
		sb.WriteString("r-xr-xr-x")
	} else {
		sb.WriteString("rwxrwxrwx")
	}

	return sb.String()
}
