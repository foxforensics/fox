package rpm

import (
	"bytes"
	"compress/bzip2"
	"io"
	"log"

	"github.com/cavaliergopher/cpio"
	"github.com/cavaliergopher/rpm"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
	"github.com/sorairolake/lzip-go"
	"github.com/ulikunitz/xz"
	"github.com/ulikunitz/xz/lzma"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte{
		0xED, 0xAB, 0xEE, 0xDB,
	})
}

func Extract(b []byte, root, _ string) (e []data.Stream) {
	br := bytes.NewReader(b)

	pkg, err := rpm.Read(br)

	if err != nil {
		log.Println(err)
		return
	}

	var r1 io.Reader

	switch v := pkg.PayloadCompression(); v {
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
		log.Printf("%s not supported\n", v)
	}

	// prevent resource leaks
	if r, ok := r1.(*gzip.Reader); ok {
		defer func() {
			_ = r.Close()
		}()
	}

	if err != nil {
		log.Println(err)
		return
	}

	if pkg.PayloadFormat() != "cpio" {
		log.Printf("%s not supported\n", pkg.PayloadFormat())
		return
	}

	r2 := cpio.NewReader(r1)

	for {
		h, err := r2.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
			break
		}

		if !h.Mode.IsRegular() {
			continue
		}

		buf := make([]byte, h.Size)

		_, err = r2.Read(buf)

		if err != nil {
			log.Println(err)
			continue
		}

		e = append(e, data.Stream{
			Path: data.JoinPart(root, h.Name),
			Data: buf,
		})
	}

	return
}
