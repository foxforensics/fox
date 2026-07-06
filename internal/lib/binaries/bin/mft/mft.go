package mft

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"

	"go.foxforensics.eu/fox/v4/internal/lib"
	"go.foxforensics.eu/fox/v4/internal/pkg/time/body"
	"www.velocidex.com/golang/go-ntfs/parser"
)

const (
	cluster = 4096 // sane size
	record  = 1024
)

func Detect(b []byte) bool {
	return lib.HasMagic(b, 0, []byte{
		'F', 'I', 'L', 'E',
	})
}

func Convert(b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	ch := parser.ParseMFTFile(context.Background(), bytes.NewReader(b), int64(len(b)), cluster, record)

	for e := range ch {
		b, err := json.Marshal(e)

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		buf.Write(b)
		buf.WriteRune('\n')
	}

	return buf.Bytes(), nil
}

func ToBody(b []byte) []body.Body {
	v := make([]body.Body, 0, len(b)/record)

	ch := parser.ParseMFTFile(context.Background(), bytes.NewReader(b), int64(len(b)), cluster, record)

	for e := range ch {
		v = append(v, body.Body{
			Name:   e.FileName(),
			Inode:  e.Inode,
			Size:   uint64(e.FileSize),
			Atime:  e.LastAccess0x30,
			Mtime:  e.LastModified0x30,
			Ctime:  e.Created0x30,
			Crtime: e.Created0x10,
		})
	}

	return v
}
