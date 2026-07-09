package journal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/Velocidex/go-journalctl/parser"
	"github.com/Velocidex/ordereddict"
	"go.foxforensics.eu/fox/v4/internal/pkg/hunt/event"
	"go.foxforensics.eu/fox/v4/library"
	"go.foxforensics.eu/fox/v4/library/binaries"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var Magic = []byte("LPKSHHRH")

var (
	ErrNoSystem    = errors.New("journal has no System section")
	ErrNoEventData = errors.New("journal has no EventData section")
)

func Detect(b []byte) bool {
	return library.HasMagic(b, 0, Magic)
}

func Convert(b []byte) ([]byte, error) {
	j, err := parser.OpenFile(bytes.NewReader(b))

	if err != nil {
		return b, err
	}

	buf := bytes.NewBuffer(nil)

	for l := range j.GetLogs(context.Background()) {
		fmt.Fprintf(buf, `{"Event": %s}`, l.String())
		buf.WriteRune('\n')
	}

	return buf.Bytes(), nil
}

func Carve(sr *io.SectionReader, off int64, cap int) <-chan *event.Event {
	ch := make(chan *event.Event, cap)

	_, err := sr.Seek(off, io.SeekStart)

	if err != nil {
		defer close(ch)
		if !errors.Is(err, io.EOF) {
			slog.Error(err.Error())
		}
		return ch
	}

	f, err := parser.OpenFile(sr)

	if err != nil {
		defer close(ch)
		slog.Error(err.Error())
		return ch
	}

	go func(f *parser.JournalFile) {
		defer close(ch)

		for evt := range f.GetLogs(context.Background()) {
			e, err := newEvent(evt)

			if err != nil {
				slog.Error(err.Error())
				continue
			}

			ch <- e
		}
	}(f)

	return ch
}

func newEvent(od *ordereddict.Dict) (*event.Event, error) {
	var sys, evt *ordereddict.Dict

	if v, ok := od.Get("System"); ok {
		if sys, ok = v.(*ordereddict.Dict); !ok {
			return nil, ErrNoSystem
		}
	} else {
		return nil, ErrNoSystem
	}

	if v, ok := od.Get("EventData"); ok {
		if evt, ok = v.(*ordereddict.Dict); !ok {
			return nil, ErrNoEventData
		}
	} else {
		return nil, ErrNoEventData
	}

	e := &event.Event{
		Host:     getString(sys, "_HOSTNAME"),
		Message:  getString(evt, "MESSAGE"),
		Source:   string(binaries.Journal),
		Category: getString(sys, "_TRANSPORT"),
		Service:  getString(sys, "_COMM"),
		Fields:   make(map[string]string),
	}

	// get timestamp in order of priority
	for _, k := range []string{
		"_SOURCE_REALTIME_TIMESTAMP",
		"SYSLOG_TIMESTAMP",
		"Timestamp",
	} {
		if v, ok := sys.Get(k); ok {
			if t, ok := v.(time.Time); ok {
				e.Time = t.UTC()
				break
			}
		}
	}

	if len(e.Message) == 0 {
		e.Message = "Undescribed event"
	}

	if v, ok := sys.GetInt64("_UID"); ok {
		e.User = fmt.Sprintf("%v", v)
	}

	if v, ok := sys.GetInt64("Seq"); ok {
		e.Sequence = fmt.Sprintf("%v", v)
	}

	if v, ok := evt.GetInt64("PRIORITY"); ok {
		if v <= 7 {
			e.Severity = 10 - int(v) // minimum 3
		} else {
			e.Severity = 10 // malformed means critical
		}
	}

	d := ordereddict.NewDict()
	d.MergeFrom(sys)
	d.MergeFrom(evt)

	for _, i := range d.Items() {
		if !strings.HasPrefix(i.Key, "(") {
			e.Fields[toTitle(i.Key)] = fmt.Sprintf("%v", i.Value)
		}
	}

	return e, nil
}

func getString(od *ordereddict.Dict, key string) string {
	v, _ := od.GetString(key)
	return v
}

func toTitle(s string) string {
	s = strings.TrimPrefix(s, "_")
	s = strings.ReplaceAll(s, "_", " ")
	s = cases.Title(language.Und).String(s)
	s = strings.ReplaceAll(s, " ", "")
	return s
}
