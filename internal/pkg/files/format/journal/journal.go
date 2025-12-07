package journal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Velocidex/go-journalctl/parser"
	"github.com/Velocidex/ordereddict"

	"github.com/cuhsat/fox/v4/internal/pkg/files"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
)

const (
	Magic = "LPKSHHRH"
)

var (
	Regex = regexp.MustCompile(Magic)
)

var (
	ErrNoSystem    = errors.New("journal has no System section")
	ErrNoEventData = errors.New("journal has no EventData section")
)

func Detect(b []byte) bool {
	return files.HasMagic(b, 0, []byte(Magic))
}

func Decode(b []byte, off int64, ext int) <-chan *event.Event {
	ch := make(chan *event.Event, 4096)

	f, err := parser.OpenFile(bytes.NewReader(b[off:]))

	if err != nil {
		defer close(ch)
		log.Println(err)
		return ch
	}

	go func(f *parser.JournalFile) {
		defer close(ch)

		for evt := range f.GetLogs(context.Background()) {
			e, r, err := newEvent(evt)

			if err != nil {
				log.Println(err)
				continue
			}

			if ext > 0 {
				addExtLevel1(e, r)
			}

			if ext > 1 {
				addExtLevel2(e, r)
			}

			ch <- e
		}
	}(f)

	return ch
}

func Convert(b []byte) ([]byte, error) {
	j, err := parser.OpenFile(bytes.NewReader(b))

	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)

	for l := range j.GetLogs(context.Background()) {
		_, err := buf.WriteString(fmt.Sprintf("%v\n", l))

		if err != nil {
			log.Println(err)
		}
	}

	return buf.Bytes(), err
}

func newEvent(od *ordereddict.Dict) (*event.Event, *ordereddict.Dict, error) {
	var sys, evt *ordereddict.Dict

	e := event.Event{
		Extension: make(map[string]any),
	}

	if v, ok := od.Get("System"); ok {
		sys = v.(*ordereddict.Dict)
	} else {
		return nil, nil, ErrNoSystem
	}

	if v, ok := od.Get("EventData"); ok {
		evt = v.(*ordereddict.Dict)
	} else {
		return nil, nil, ErrNoEventData
	}

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

	e.Host, _ = sys.GetString("_HOSTNAME")
	e.User, _ = sys.GetString("_UID")
	e.Message, _ = evt.GetString("MESSAGE")

	if len(e.Message) == 0 {
		e.Message = "Undescribed event"
	}

	if v, ok := evt.GetInt64("PRIORITY"); ok {
		e.Severity = 10 - int8(v) // minimum 3
	}

	r := ordereddict.NewDict()
	r.MergeFrom(sys)
	r.MergeFrom(evt)

	return &e, r, nil
}

func addExtLevel1(e *event.Event, od *ordereddict.Dict) {
	e.Extension["rt"] = e.Time
	e.Extension["shost"] = e.Host
	e.Extension["suid"] = e.User
	e.Extension["deviceFacility"] = "systemd"

	for k, v := range map[string]string{
		"cat":               "_TRANSPORT",
		"spid":              "_PID",
		"sourceServiceName": "_COMM",
	} {
		if a, ok := od.Get(v); ok {
			e.Extension[k] = fmt.Sprintf("%v", a)
		}
	}

	return
}

func addExtLevel2(e *event.Event, od *ordereddict.Dict) {
	for _, i := range od.Items() {
		if !strings.HasPrefix(i.Key, "(") {
			e.Extension[i.Key] = fmt.Sprintf("%v", i.Value)
		}
	}
}
