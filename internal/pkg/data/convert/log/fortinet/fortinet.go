// Package fortinet source:
// https://cyber.wtf/2024/08/30/parsing-fortinet-binary-firewall-logs/
// https://github.com/GDATAAdvancedAnalytics/FortilogDecoder/blob/main/fortilog_decoder.py
package fortinet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/pierrec/lz4/v4"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

var (
	ErrTlcDeflate   = errors.New("can not deflate TLC")
	ErrNotSupported = errors.New("type not supported")
)

var magic = [][]byte{
	{0xEC, 0xCF},
	{0xEC, 0xDE},
}

// Timestamp     uint32 TODO

func Detect(b []byte) bool {
	for _, m := range [][]byte{
		{0xEC, 0xCE}, // llog v5 old
		{0xEC, 0xCF}, // llog v5 old
		{0xEC, 0xDF}, // llog v5 new
		{0xAA, 0x01}, // tlc
	} {
		if data.HasMagic(b, 0, m) {
			return true
		}
	}

	return false
}

func Convert(b []byte) ([]byte, error) {
	var buf bytes.Buffer

	err := decodeLLogV5(b, &buf)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decodeLLogV5(b []byte, out *bytes.Buffer) error {
	r := bytes.NewReader(b)

	var i int64

	for {
		i = r.Size() - int64(r.Len())

		// Peek at next 2 bytes to determine type
		logType := make([]byte, 2)
		n, err := r.Read(logType)
		if err == io.EOF || n < 2 {
			break
		}
		if err != nil {
			return err
		}
		_, _ = r.Seek(-2, io.SeekCurrent) // Unread for processing

		if bytes.Equal(logType, magic[0]) || bytes.Equal(logType, magic[1]) {
			// Consume magic bytes
			_, _ = r.Read(logType)

			head := make([]byte, 16)
			if _, err := io.ReadFull(r, head); err != nil {
				return err
			}

			flag := (head[0] >> 2) & 1
			lDevID := int(head[3])
			lDevName := int(head[4])
			lVDOM := int(head[5])
			entryCount := int(binary.BigEndian.Uint16(head[6:8]))
			lCompressed := int(binary.BigEndian.Uint16(head[8:10]))
			lDecompressed := int(binary.BigEndian.Uint16(head[10:12]))
			lEntryCounts := entryCount * 2
			lSomething := 0

			if flag != 0 {
				lSomething = lEntryCounts
			}

			var tzLen int
			if bytes.Equal(logType, magic[1]) { // 0xECDE
				_, _ = r.Seek(10, io.SeekCurrent)
				tzByte, _ := r.ReadByte()
				tzLen = int(tzByte)
			}

			lASCII := lDevID + lDevName + lVDOM
			body := make([]byte, lASCII+lEntryCounts+lSomething)

			if _, err := io.ReadFull(r, body); err != nil {
				return err
			}

			devID := string(body[0:lDevID])
			devName := string(body[lDevID : lDevID+lDevName])
			vdom := string(body[lDevID+lDevName : lDevID+lDevName+lVDOM])
			entriesLengths := body[lASCII : lASCII+lEntryCounts]

			if bytes.Equal(logType, magic[1]) && tzLen > 0 {
				_, _ = r.Seek(int64(tzLen), io.SeekCurrent)
			}

			compressed := make([]byte, lCompressed)
			if _, err := io.ReadFull(r, compressed); err != nil {
				return err
			}

			decompressed := make([]byte, lDecompressed+1)
			uncomp, err := lz4.UncompressBlock(compressed, decompressed)
			if err != nil {
				// Skip this entry on decompression error, continue processing
				continue
			}
			decompressed = decompressed[:uncomp]

			prefix := fmt.Sprintf(`devid="%s" devname="%s" vdom="%s" `, devID, devName, vdom)

			if entryCount > 1 {
				pointer := 0
				for i := 0; i < lEntryCounts; i += 2 {
					l := int(binary.BigEndian.Uint16(entriesLengths[i : i+2]))
					out.WriteString(prefix)
					out.Write(decompressed[pointer : pointer+l])
					out.WriteByte(0x0a)
					pointer += l
				}
			} else if entryCount == 1 {
				out.WriteString(prefix)
				out.Write(decompressed)
				out.WriteByte(0x0a)
			}

			// Skip 2nd variable part
			head2 := make([]byte, 2)
			if _, err := io.ReadFull(r, head2); err != nil {
				continue
			}
			body2 := binary.LittleEndian.Uint16(head2)
			_, _ = r.Seek(int64(body2), io.SeekCurrent)

		} else if bytes.Equal(logType, []byte{0xAA, 0x01}) {
			_, _ = r.Seek(4, io.SeekCurrent) // Skip magic + 2 bytes

			lBodyBytes := make([]byte, 4)
			if _, err := io.ReadFull(r, lBodyBytes); err != nil {
				return err
			}
			lBody := int(binary.BigEndian.Uint32(lBodyBytes)) - 8

			body := make([]byte, lBody)
			if _, err := io.ReadFull(r, body); err != nil {
				return err
			}

			tlc, err := parseTLC(body)
			if err != nil {
				continue
			}

			rawEntries := bytes.Split(tlc, []byte{0x00})
			for _, rawEntry := range rawEntries {
				idx := bytes.Index(rawEntry, []byte("date="))
				if idx == -1 {
					continue
				}
				out.Write(rawEntry[idx:])
				out.WriteByte(0x0a)
			}

		} else if bytes.Equal(logType, []byte{0x00, 0x00}) || logType[0] == 0x00 {
			_, _ = r.ReadByte()
			continue
		} else {
			return fmt.Errorf("unknown header %x at offset %d", logType, i)
		}
	}

	return nil
}

func parseTLC(b []byte) ([]byte, error) {
	var zBufLen int
	var zBuf []byte

	for i := 0; i < len(b); {
		var l int
		var v int64

		if i >= len(b) {
			break
		}
		typeHigh := b[i] >> 4
		i++

		if i >= len(b) {
			break
		}
		field := b[i]
		i++

		switch typeHigh {
		case 0: // byte array prefixed with int8 length
			i, l = i+1, int(b[i])

		case 1: // byte array prefixed with int16be length
			i, l = i+2, int(binary.BigEndian.Uint16(b[i:i+2]))

		case 2: // byte array prefixed with int32be length
			i, l = i+4, int(binary.BigEndian.Uint32(b[i:i+4]))

		case 3: // int8
			i, v = i+1, int64(b[i])

		case 4: // int16be
			i, v = i+2, int64(binary.BigEndian.Uint16(b[i:i+2]))

		case 5: // int32be
			i, v = i+4, int64(binary.BigEndian.Uint32(b[i:i+4]))

		case 6: // int64be
			i, v = i+8, int64(binary.BigEndian.Uint64(b[i:i+8]))

		default:
			return nil, ErrNotSupported
		}

		// type is byte array
		if typeHigh <= 2 {
			i, zBuf = i+l, b[i:i+l]
		}

		switch field {
		case 12: // set buf len
			zBufLen = int(v)

		case 16: // deflate buf
			if len(zBuf) == 0 || zBufLen == 0 {
				return nil, ErrTlcDeflate
			}

			r := lz4.NewReader(bytes.NewReader(zBuf))
			d := make([]byte, zBufLen)
			n, err := io.ReadFull(r, d)

			if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
				return nil, err
			}

			return d[:n], nil
		}
	}

	return nil, ErrTlcDeflate
}
