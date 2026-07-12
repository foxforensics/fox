package evtx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Velocidex/ordereddict"
	"go.foxforensics.eu/eventid/events"
	"go.foxforensics.eu/fox/v5/internal/cmd/hunt/event"
	"go.foxforensics.eu/fox/v5/library"
	"go.foxforensics.eu/fox/v5/library/binaries"
	"www.velocidex.com/golang/evtx"
)

var Magic = []byte(evtx.EVTX_HEADER_MAGIC)
var Chunk = []byte(evtx.EVTX_CHUNK_HEADER_MAGIC)

var providers events.Providers

func Detect(b []byte) bool {
	return library.HasMagic(b, 0, Magic)
}

func Convert(b []byte) ([]byte, error) {
	chunks, err := evtx.GetChunks(bytes.NewReader(b))

	if err != nil {
		return b, err
	}

	buf := bytes.NewBuffer(nil)

	for _, chunk := range chunks {
		records, err := chunk.Parse(0)

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		for _, record := range records {
			b, err := json.Marshal(record.Event)

			if err != nil {
				slog.Error(err.Error())
				continue
			}

			buf.Write(b)
			buf.WriteRune('\n')
		}
	}

	return buf.Bytes(), nil
}

func Prepare() {
	var err error

	// use empty fallback on errors
	if providers, err = events.Get(); err != nil {
		providers = make(events.Providers)
		slog.Error(err.Error())
	}
}

func Carve(sr *io.SectionReader, off int64, cap int) <-chan *event.Event {
	ch := make(chan *event.Event, cap)

	chunk, err := evtx.NewChunk(sr, off)

	if err != nil {
		defer close(ch)

		if !errors.Is(err, io.EOF) {
			slog.Error(err.Error())
		}

		return ch
	}

	go func() {
		defer close(ch)
		records, err := chunk.Parse(0)

		if err != nil {
			slog.Error(err.Error())
			return
		}

		for _, record := range records {
			e, err := newEvent(record)

			if err != nil {
				slog.Error(err.Error())
				continue
			}

			ch <- e
		}
	}()

	return ch
}

func newEvent(r *evtx.EventRecord) (*event.Event, error) {
	od, ok := r.Event.(*ordereddict.Dict)

	if !ok {
		return nil, errors.New("event type invalid")
	}

	e := &event.Event{
		Time:     intToUTC(r.Header.FileTime),
		Host:     getString(od, "Event.System.Computer"),
		User:     getAny(od, "Event.System.Security.UserID"),
		Sequence: strconv.Itoa(getInt(od, "Event.System.EventRecordID")),
		Source:   string(binaries.Eventlog),
		Category: getString(od, "Event.System.Channel"),
		Service:  getString(od, "Event.System.Provider.Name"),
		Fields:   make(map[string]string),
	}

	// fallback service name
	if len(e.Service) == 0 {
		e.Service = "unknown"
	}

	e.Message = fmt.Sprintf("Undescribed event: %s: %d", e.Service, r.Header.RecordID)

	// calculate severity from level
	switch getInt(od, "Event.System.Level") {
	case 0, 1:
		e.Severity = 10
	case 2:
		e.Severity = 8
	case 3:
		e.Severity = 6
	case 4:
		e.Severity = 3
	case 5:
		e.Severity = 1
	default:
		e.Severity = 0
	}

	addFields(e, od, "")

	// translate event id to message
	if provider, ok := providers[e.Service]; ok {
		if id := int64(getInt(od, "Event.System.EventID.Value")); id > 0 {
			if message, ok := provider[id]; ok {
				e.Message = expandParams(message, e)
			}
		}
	}

	return e, nil
}

func addFields(e *event.Event, od *ordereddict.Dict, parent string) {
	if omitParent(parent) {
		parent = ""
	}

	for _, item := range od.Items() {
		var key = parent
		var val string

		switch v := item.Value.(type) {
		case *ordereddict.Dict:
			addFields(e, v, item.Key)
			continue

		case float64: // is always a filetime
			val = floatToUTC(v).Format(time.RFC3339Nano)

		default:
			val = fmt.Sprintf("%v", v)
		}

		if od.Len() > 1 || !omitChild(item.Key) {
			key += item.Key // combined key
		}

		if len(key) == 0 {
			key = item.Key // safety
		}

		e.Fields[key] = val
	}
}

func expandParams(msg string, e *event.Event) string {
	if strings.Contains(msg, "$1") {
		for i := 1; ; i++ {
			if v, ok := e.Fields[fmt.Sprintf("param%d", i)]; ok {
				msg = strings.ReplaceAll(msg, fmt.Sprintf("$%d", i), v)
			} else {
				break // no more params
			}
		}
	}

	return msg
}

func omitParent(key string) bool {
	return slices.Contains([]string{
		"System",
		"Security",
		"Execution",
		"EventData",
	}, key)
}

func omitChild(key string) bool {
	return slices.Contains([]string{
		"Guid",
		"Value",
		"SystemTime",
	}, key)
}

func getString(od *ordereddict.Dict, key string) string {
	v, ok := ordereddict.GetString(od, key)
	if !ok {
		slog.Error(fmt.Sprintf("%s is not a string", key))
	}
	return v
}

func getInt(od *ordereddict.Dict, key string) int {
	v, ok := ordereddict.GetInt(od, key)
	if !ok {
		slog.Error(fmt.Sprintf("%s is not an int", key))
	}
	return v
}

func getAny(od *ordereddict.Dict, key string) string {
	v, ok := ordereddict.GetAny(od, key)
	if !ok {
		return "" // fallback for nil
	}
	return fmt.Sprintf("%v", v)
}

func floatToUTC(v float64) time.Time {
	nsec := int64((v - float64(int64(v))) * 1e9)
	return time.Unix(int64(v), nsec).UTC()
}

func intToUTC(v uint64) time.Time {
	const nsec int64 = 116444736000000000
	return time.Unix(0, (int64(v)-nsec)*100).UTC()
}
