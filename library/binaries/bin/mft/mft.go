package mft

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"go.foxforensics.eu/fox/v5/internal/cmd/time/entry"
	"go.foxforensics.eu/fox/v5/library"
	"www.velocidex.com/golang/go-ntfs/parser"
)

const (
	cluster = 4096 // sane size
	record  = 1024
)

func Detect(b []byte) bool {
	return library.HasMagic(b, 0, []byte{
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

func Parse(b []byte) []entry.Entry {
	v := make([]entry.Entry, 0, len(b)/record)

	ch := parser.ParseMFTFile(context.Background(), bytes.NewReader(b), int64(len(b)), cluster, record)

	for mh := range ch {
		v = append(v, entry.Entry{
			Name:    mh.FileName(),
			Inode:   fmt.Sprintf("%d-%d", mh.EntryNumber, mh.SequenceNumber),
			Mode:    buildMode(mh),
			Size:    uint64(mh.FileSize),
			Mtime:   mh.LastModified0x10.UTC(),
			Atime:   mh.LastAccess0x10.UTC(),
			Ctime:   mh.LastRecordChange0x10.UTC(),
			Btime:   mh.Created0x10.UTC(),
			Anomaly: checkTimes(mh),
		})
	}

	return v
}

func buildMode(mh *parser.MFTHighlight) string {
	var sb strings.Builder

	if strings.Contains(mh.SIFlags, "REPARSE_POINT") {
		sb.WriteByte('l')
	} else if mh.IsDir {
		sb.WriteByte('d')
	} else {
		sb.WriteByte('-')
	}

	if strings.Contains(mh.SIFlags, "READ_ONLY") || strings.Contains(mh.SIFlags, "SYSTEM") {
		sb.WriteString("r-xr-xr-x")
	} else {
		sb.WriteString("rwxrwxrwx")
	}

	return sb.String()
}

func checkTimes(mh *parser.MFTHighlight) bool {
	return checkTime(mh.LastModified0x10, mh.LastModified0x30) ||
		checkTime(mh.LastAccess0x10, mh.LastAccess0x30) ||
		checkTime(mh.LastRecordChange0x10, mh.LastRecordChange0x30) ||
		checkTime(mh.Created0x10, mh.Created0x30)
}

func checkTime(t1, t2 time.Time) bool {
	if t1.Nanosecond() == 0 {
		if t1.UTC().Unix() <= t2.UTC().Unix() {
			return true // with high probability
		}
	}

	return false
}
