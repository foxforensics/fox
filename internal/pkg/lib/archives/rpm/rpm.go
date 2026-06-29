package rpm

import (
	"bytes"
	"compress/bzip2"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/cavaliergopher/cpio"
	"github.com/cavaliergopher/rpm"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
	"github.com/sorairolake/lzip-go"
	"github.com/ulikunitz/xz"
	"github.com/ulikunitz/xz/lzma"
	"go.foxforensics.eu/fox/v4/internal/pkg/lib"
	"go.foxforensics.eu/fox/v4/internal/sys"
)

func Detect(b []byte) bool {
	return lib.HasMagic(b, 0, []byte{
		0xED, 0xAB, 0xEE, 0xDB,
	})
}

func Extract(b []byte, root, _ string) (e []lib.Stream) {
	br := bytes.NewReader(b)

	rp, err := rpm.Read(br)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	var r1 io.Reader

	switch v := rp.PayloadCompression(); v {
	case "bzip2":
		r1 = bzip2.NewReader(br)
	case "gzip":
		r1, err = gzip.NewReader(br)
	case "lzip":
		r1, err = lzip.NewReader(br)
	case "lzma":
		r1, err = lzma.NewReader(br)
	case "xz":
		r1, err = xz.NewReader(br)
	case "zstd":
		r1, err = zstd.NewReader(br)
	case "uncompressed":
		r1 = br
	default:
		slog.Error(fmt.Sprintf("%s not supported!", v))
		return
	}

	if err != nil {
		slog.Error(err.Error())
		return
	}

	// prevent resource leaks
	switch r1.(type) {
	case *gzip.Reader:
		defer func() {
			if err := r1.(*gzip.Reader).Close(); err != nil {
				slog.Error(err.Error())
			}
		}()

	case *zstd.Decoder:
		defer r1.(*zstd.Decoder).Close()
	}

	if rp.PayloadFormat() != "cpio" {
		slog.Warn(fmt.Sprintf("%s not supported!", rp.PayloadFormat()))
		return
	}

	r2 := cpio.NewReader(r1)

	for {
		h, err := r2.Next()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			slog.Error(err.Error())
			break
		}

		if !h.Mode.IsRegular() {
			continue
		}

		buf := make([]byte, h.Size)

		_, err = io.ReadFull(r2, buf)

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		e = append(e, lib.Stream{
			Path: sys.JoinPart(root, h.Name),
			Data: buf,
		})
	}

	return
}
