package fortinet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/pierrec/lz4/v4"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

type llog5 struct {
	// entries raw
	Magic         uint16
	Flags         uint8
	Unused        uint16
	LDevId        uint8
	LDevName      uint8
	LVDom         uint8
	Entries       uint16
	LCompressed   uint16
	LDecompressed uint16
	Timestamp     uint32

	// entries parsed
	LEntries uint16
	LAscii   uint16
	Padding  uint16
	DevId    string
	DevName  string
	VDom     string
	Body     []byte
}

func Detect(b []byte) bool {
	for _, m := range [][]byte{
		{0xEC, 0xCE}, // llog v5
		{0xEC, 0xCF}, // llog v5
	} {
		if data.HasMagic(b, 0, m) {
			return true
		}
	}

	return false
}

func Convert(b []byte) ([]byte, error) {
	log.Println("warning: fortinet parser is experimental!")

	buf := bytes.NewBuffer(nil)

	r := bytes.NewReader(b)

	for {
		l, err := decode(r)

		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}

		d, err := deflate(l)

		if err != nil {
			return nil, err
		}

		buf.Write(d)

		err = forward(r)

		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func decode(r io.Reader) (*llog5, error) {
	var err error

	llog := new(llog5)

	err = binary.Read(r, binary.LittleEndian, &llog.Magic)

	if llog.Magic != 0xCEEC && llog.Magic != 0xCFEC {
		return nil, errors.New("log format not supported")
	}

	err = binary.Read(r, binary.LittleEndian, &llog.Flags)
	err = binary.Read(r, binary.LittleEndian, &llog.Unused)
	err = binary.Read(r, binary.LittleEndian, &llog.LDevId)
	err = binary.Read(r, binary.LittleEndian, &llog.LDevName)
	err = binary.Read(r, binary.LittleEndian, &llog.LVDom)
	err = binary.Read(r, binary.BigEndian, &llog.Entries)
	err = binary.Read(r, binary.BigEndian, &llog.LCompressed)
	err = binary.Read(r, binary.BigEndian, &llog.LDecompressed)
	err = binary.Read(r, binary.BigEndian, &llog.Timestamp)

	if err != nil {
		return llog, err
	}

	llog.LEntries = llog.Entries * 2
	llog.LAscii = uint16(llog.LDevId + llog.LDevName + llog.LVDom)

	if llog.Flags&4 == 1 {
		llog.Padding = llog.LEntries
	}

	llog.Body = make([]byte, llog.LAscii+llog.LEntries+llog.Padding+llog.LCompressed)

	_, _ = io.ReadFull(r, llog.Body)

	i, j := 0, int(llog.LDevId)

	llog.DevId = string(llog.Body[i:j])

	i, j = j, j+int(llog.LDevName)

	llog.DevName = string(llog.Body[i:j])

	i, j = j, j+int(llog.LVDom)

	llog.VDom = string(llog.Body[i:j])

	return llog, nil
}

func deflate(llog *llog5) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	i := llog.LAscii + llog.LEntries + llog.Padding
	j := i + llog.LCompressed

	b := make([]byte, llog.LDecompressed+1)

	n, err := lz4.UncompressBlock(llog.Body[i:j], b)

	if err != nil {
		return nil, err
	}

	if uint16(n) != llog.LDecompressed {
		return nil, errors.New("invalid block length")
	}

	if llog.Entries == 1 {
		_, _ = buf.WriteString(format(llog, b))
	} else {
		p, q, b := 0, 0, llog.Body[llog.LAscii:llog.LAscii+llog.LEntries]

		for i := 0; i < int(llog.LEntries); i += 2 {
			q = int(binary.BigEndian.Uint16(b[i : i+2]))
			_, _ = buf.WriteString(format(llog, b[p:p+q]))
			p += q
		}
	}

	return buf.Bytes(), nil
}

func format(llog *llog5, b []byte) string {
	return fmt.Sprintf("devid=\"%s\" devname=\"%s\" vdom=\"%s\" %s\n", llog.DevId, llog.DevName, llog.VDom, string(b))
}

func forward(r io.ReadSeeker) error {
	var n uint16

	err := binary.Read(r, binary.LittleEndian, &n)

	if err != nil {
		return err
	}

	_, err = r.Seek(int64(n), io.SeekCurrent)

	return err
}
