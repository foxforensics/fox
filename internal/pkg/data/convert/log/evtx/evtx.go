package evtx

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"maps"
	"slices"
	"strconv"
	"time"

	"github.com/0xrawsec/golang-evtx/evtx"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
)

var Magic = []byte(evtx.EvtxMagic)
var Chunk = []byte(evtx.ChunkMagic)

var system = evtx.Path("/Event/System")

var children = []string{
	"Guid",
	"Name",
	"Value",
	"DwordVal",
}

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, Magic)
}

func Convert(b []byte) ([]byte, error) {
	r, err := evtx.New(bytes.NewReader(b))

	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)

	for e := range r.Events() {
		buf.Write(evtx.ToJSON(e))
		buf.WriteRune('\n')
	}

	_ = r.Close()

	return buf.Bytes(), nil
}

func Carve(rs io.ReadSeeker, off int, cap int) <-chan *event.Event {
	ch := make(chan *event.Event, cap)

	chk, err := newChunk(rs, int64(off))

	if err != nil {
		defer close(ch)
		log.Println(err)
		return ch
	}

	go func(chk *evtx.Chunk) {
		defer close(ch)
		for evt := range chk.Events() {
			ch <- newEvent(evt)
		}
	}(chk)

	return ch
}

func newChunk(rs io.ReadSeeker, off int64) (*evtx.Chunk, error) {
	evtx.SetModeCarving(true)
	evtx.GoToSeeker(rs, off)

	chk := evtx.NewChunk()
	chk.Offset = off
	chk.Data = make([]byte, evtx.ChunkSize)

	if _, err := rs.Read(chk.Data); err != nil {
		return nil, err
	}

	r := bytes.NewReader(chk.Data)

	chk.ParseChunkHeader(r)

	if err := chk.Header.Validate(); err != nil {
		return nil, err
	}

	evtx.GoToSeeker(r, int64(chk.Header.SizeHeader))

	chk.ParseStringTable(r)

	if err := chk.ParseTemplateTable(r); err != nil {
		return nil, err
	}

	if err := chk.ParseEventOffsets(r); err != nil {
		return nil, err
	}

	return &chk, nil
}

func newEvent(evt *evtx.GoEvtxMap) *event.Event {
	e := event.Event{
		Time:     evt.TimeCreated().UTC(),
		Host:     getString(evt, "/Event/System/Computer"),
		User:     getString(evt, "/Event/System/Security/UserID"),
		Sequence: getString(evt, "/Event/System/EventRecordID"),
		Source:   types.Eventlog,
		Category: getString(evt, "/Event/System/Channel"),
		Service:  getString(evt, "/Event/System/Provider/Name"),
		Fields:   make(map[string]string),
	}

	// fallback service name
	if len(e.Service) == 0 {
		e.Service = "unknown"
	}

	// translate event id to message
	e.Message = fmt.Sprintf("Undescribed event: %s: %d", e.Service, evt.EventID())

	if events, ok := Events[e.Service]; ok {
		if message, ok := events[evt.EventID()]; ok {
			e.Message = message
		}
	}

	// calculate severity from level
	level := getString(evt, "/Event/System/Level")

	if v, err := strconv.Atoi(level); err == nil {
		switch v {
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
	}

	if v, err := evt.GetMap(&system); err == nil {
		addMapDeep(&e, v, "")
	}

	return &e
}

func addMapDeep(e *event.Event, em *evtx.GoEvtxMap, p string) {
	for k, v := range maps.All(*em) {
		switch v.(type) {
		case *evtx.GoEvtxMap, evtx.GoEvtxMap:
			m := v.(evtx.GoEvtxMap)
			addMapDeep(e, &m, k)

		case *evtx.UTCTime, evtx.UTCTime:
			u := v.(evtx.UTCTime)
			e.Fields[k] = fmt.Sprintf("%s", time.Time(u))

		default:
			if slices.Contains(children, k) {
				k = p // use parent as key
			}

			e.Fields[k] = fmt.Sprintf("%v", v)
		}
	}
}

func getString(em *evtx.GoEvtxMap, path string) string {
	p := evtx.Path(path)
	v, _ := em.GetString(&p)
	return v
}
