// Package ecs specification:
// https://www.elastic.co/docs/reference/ecs/ecs-field-reference
package ecs

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cuhsat/fox/v4/internal"
	"github.com/cuhsat/fox/v4/internal/pkg/data/stream"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
)

const version = "9.1.0"

type Ecs struct {
	stream.Stream

	Ecs struct {
		Version string `json:"version"`
	} `json:"ecs"`

	Agent struct {
		Type    string `json:"type"`
		Version string `json:"version"`
	} `json:"agent"`

	Host struct {
		Hostname string `json:"hostname,omitempty"`
	} `json:"host,omitempty"`

	User struct {
		ID string `json:"id,omitempty"`
	} `json:"user,omitempty"`

	Event struct {
		Kind     string    `json:"kind,omitempty"`
		Module   string    `json:"module,omitempty"`
		Dataset  string    `json:"dataset,omitempty"`
		Severity int64     `json:"severity,omitempty"`
		ID       string    `json:"id,omitempty"`
		Code     string    `json:"code,omitempty"`
		Provider string    `json:"provider,omitempty"`
		Ingested time.Time `json:"ingested,omitempty"`
		Original string    `json:"original,omitempty"`
		Hash     string    `json:"hash,omitempty"`
	} `json:"event"`

	Labels map[string]any `json:"labels,omitempty"`

	Timestamp time.Time `json:"@timestamp"`
	Message   string    `json:"message"`
}

func New(url string) Ecs {
	ecs := Ecs{Stream: stream.Stream{Url: url, Map: map[string]string{
		"Content-Type": "application/json",
	}}}

	ecs.Ecs.Version = version
	ecs.Agent.Type = "fox"
	ecs.Agent.Version = app.Version[1:]
	ecs.Event.Kind = "event"

	return ecs
}

func (ecs Ecs) String() string {
	return fmt.Sprintf("ECS: %s", ecs.Url)
}

func (ecs Ecs) Write(e *event.Event) (int64, int64, error) {
	cef := e.ToCEF()

	ecs.Timestamp = e.Time.UTC()
	ecs.Message = e.Message

	ecs.Host.Hostname = e.Host
	ecs.User.ID = e.User

	switch e.Source {

	// windows specific
	case types.Eventlog:
		ecs.Event.Module = "eventlog"
		ecs.Event.Dataset = "EventLog." + e.Value("System_Channel")
		ecs.Event.ID = e.Value("System_EventRecordID")
		ecs.Event.Code = e.Value("System_EventID", "System_EventID_Value")
		ecs.Event.Provider = e.Value("System_Provider_Name")

	// linux specific
	case types.Journal:
		ecs.Event.Module = "journal"
		ecs.Event.Dataset = e.Value("_COMM")
		ecs.Event.ID = e.Value("Seq")
		ecs.Event.Provider = e.Value("_TRANSPORT")
	}

	// original event
	ecs.Event.Severity = int64(e.Severity)
	ecs.Event.Ingested = time.Now().UTC()
	ecs.Event.Original = cef
	ecs.Event.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(cef)))

	ecs.Labels = make(map[string]any)

	for k, v := range e.Extension {
		ecs.Labels[k] = v
	}

	buf, err := json.Marshal(ecs)

	if err != nil {
		return 0, 0, nil
	}

	return ecs.Post(string(buf))
}
