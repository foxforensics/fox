package iso

import (
	"bytes"
	"io"
	"log"
	"path/filepath"

	"github.com/hooklift/iso9660"

	"github.com/cuhsat/fox/v4/internal/pkg/file"
)

func Detect(b []byte) bool {
	for _, o := range []int{
		0x8001, 0x8801, 0x9001,
	} {
		if file.HasMagic(b, o, []byte{
			'C', 'D', '0', '0', '1',
		}) {
			return true
		}
	}

	return false
}

func Extract(b []byte, root, _ string) (e []file.Stream) {
	r, err := iso9660.NewReader(bytes.NewReader(b))

	if err != nil {
		log.Println(err)
		return
	}

	for {
		f, err := r.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
			break
		}

		if f.IsDir() {
			continue
		}

		buf, err := io.ReadAll(f.Sys().(io.Reader))

		if err != nil {
			log.Println(err)
			continue
		}

		e = append(e, file.Stream{
			Path: file.JoinPart(root, filepath.Base(f.Name())),
			Data: buf,
		})
	}

	return
}
