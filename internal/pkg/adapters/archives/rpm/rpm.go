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
	"go.foxforensics.eu/fox/v4/internal/pkg"
	"go.foxforensics.eu/fox/v4/internal/sys"
)

func Detect(b []byte) bool {
	return pkg.HasMagic(b, 0, []byte{
		0xED, 0xAB, 0xEE, 0xDB,
	})
}

func Extract(b []byte, root, _ string) (e []pkg.Stream) {
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
		slog.Warn(fmt.Sprintf("%s not supported!", v))
	}

	// prevent resource leaks
	if r, ok := r1.(*gzip.Reader); ok {
		defer func() {
			if err = r.Close(); err != nil {
				slog.Error(err.Error())
			}
		}()
	}

	if err != nil {
		slog.Error(err.Error())
		return
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

		_, err = r2.Read(buf)

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		e = append(e, pkg.Stream{
			Path: sys.JoinPart(root, h.Name),
			Data: buf,
		})
	}

	return
}
