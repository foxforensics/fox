package journal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/Velocidex/go-journalctl/parser"
	"github.com/Velocidex/ordereddict"

	"foxhunt.dev/fox/internal/pkg/data"
	"foxhunt.dev/fox/internal/pkg/types"
	"foxhunt.dev/fox/internal/pkg/types/event"
)

var Magic = []byte("LPKSHHRH")

var (
	ErrNoSystem    = errors.New("journal has no System section")
	ErrNoEventData = errors.New("journal has no EventData section")
)

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, Magic)
}

func Convert(b []byte) ([]byte, error) {
	j, err := parser.OpenFile(bytes.NewReader(b))

	if err != nil {
		return b, err
	}

	buf := bytes.NewBuffer(nil)

	for l := range j.GetLogs(context.Background()) {
		buf.WriteString(fmt.Sprintf(`{"Event": %s}`, l.String()))
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
			log.Println(err)
		}
		return ch
	}

	f, err := parser.OpenFile(sr)

	if err != nil {
		defer close(ch)
		log.Println(err)
		return ch
	}

	go func(f *parser.JournalFile) {
		defer close(ch)

		for evt := range f.GetLogs(context.Background()) {
			e, err := newEvent(evt)

			if err != nil {
				log.Println(err)
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
		sys = v.(*ordereddict.Dict)
	} else {
		return nil, ErrNoSystem
	}

	if v, ok := od.Get("EventData"); ok {
		evt = v.(*ordereddict.Dict)
	} else {
		return nil, ErrNoEventData
	}

	e := event.Event{
		Host:     getString(sys, "_HOSTNAME"),
		Message:  getString(evt, "MESSAGE"),
		Source:   types.Journal,
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
			e.Time = v.(time.Time).UTC()
			break
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
		e.Severity = 10 - int(v) // minimum 3
	}

	d := ordereddict.NewDict()
	d.MergeFrom(sys)
	d.MergeFrom(evt)

	for _, i := range d.Items() {
		if !strings.HasPrefix(i.Key, "(") {
			e.Fields[i.Key] = fmt.Sprintf("%v", i.Value)
		}
	}

	return &e, nil
}

func getString(od *ordereddict.Dict, path string) string {
	v, _ := od.GetString(path)
	return v
}
