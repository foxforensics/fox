package evtx

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"maps"
	"regexp"
	"time"

	"github.com/0xrawsec/golang-evtx/evtx"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
)

const (
	Magic = evtx.EvtxMagic
	Chunk = evtx.ChunkMagic
)

var Regex = regexp.MustCompile(Chunk)

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte(Magic))
}

func Format(b []byte, _ int) ([]byte, error) {
	r, err := evtx.New(bytes.NewReader(b))

	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)

	for e := range r.Events() {
		_, err := buf.Write(evtx.ToJSON(e))

		if err != nil {
			log.Println(err)
			continue
		}

		_, err = buf.WriteRune('\n')

		if err != nil {
			log.Println(err)
		}
	}

	_ = r.Close()

	return buf.Bytes(), nil
}

func Parse(rs io.ReadSeeker, off int64, ext int) <-chan *event.Event {
	ch := make(chan *event.Event, types.Buffer)

	chk, err := newChunk(rs, off)

	if err != nil {
		defer close(ch)
		log.Println(err)
		return ch
	}

	go func(chk *evtx.Chunk) {
		defer close(ch)

		for evt := range chk.Events() {
			e := newEvent(evt)

			if ext > 0 {
				addExtLevel1(e, evt)
			}

			if ext > 1 {
				addExtLevel2(e, evt)
			}

			ch <- e
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
	var ok bool

	e := event.Event{
		Time:      evt.TimeCreated().UTC(),
		Source:    types.Eventlog,
		Extension: make(map[string]any),
	}

	hostPath := evtx.Path("/Event/System/Computer")
	namePath := evtx.Path("/Event/System/Provider/Name")

	provider, _ := evt.GetString(&namePath)
	e.Host, _ = evt.GetString(&hostPath)
	e.User, _ = evt.UserID()

	if len(provider) == 0 {
		provider = "unknown"
	}

	e.Message = fmt.Sprintf("Undescribed event: %s: %d", provider, evt.EventID())

	if events, ok := Events[provider]; ok {
		if message, ok := events[evt.EventID()]; ok {
			e.Message = message
		}
	}

	if e.Severity, ok = Levels[evt.EventID()]; !ok {
		e.Severity = 0 // unknown
	}

	return &e
}

func addExtLevel1(e *event.Event, em *evtx.GoEvtxMap) {
	e.Extension["rt"] = e.Time
	e.Extension["shost"] = e.Host
	e.Extension["suid"] = e.User
	e.Extension["deviceFacility"] = "eventlog"

	for k, v := range map[string]string{
		"cat":               "/Event/System/Channel",
		"spid":              "/Event/System/Execution/ProcessID",
		"sourceServiceName": "/Event/System/Provider/Name",
	} {
		p := evtx.Path(v)
		if s, err := em.GetString(&p); err == nil {
			e.Extension[k] = s
		}
	}
}

func addExtLevel2(e *event.Event, em *evtx.GoEvtxMap) {
	p := evtx.Path("/Event/System")

	evt, err := em.GetMap(&p)

	if err == nil {
		addMapDeep(e, evt, "")
	}
}

func addMapDeep(e *event.Event, em *evtx.GoEvtxMap, r string) {
	if len(r) > 0 {
		r += "_"
	}

	for k, v := range maps.All(*em) {
		switch v.(type) {
		case *evtx.GoEvtxMap, evtx.GoEvtxMap:
			m := v.(evtx.GoEvtxMap)
			addMapDeep(e, &m, r+k)

		case *evtx.UTCTime, evtx.UTCTime:
			u := v.(evtx.UTCTime)
			e.Extension[r+k] = fmt.Sprintf("%s", time.Time(u))

		default:
			e.Extension[r+k] = fmt.Sprintf("%v", v)
		}
	}
}
