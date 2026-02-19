// Package ecs specification:
// https://www.elastic.co/docs/reference/ecs/ecs-field-reference
package ecs

import (
	"encoding/json"
	"fmt"
	"time"

	"foxhunt.dev/fox/internal/pkg/types"
	"github.com/twmb/murmur3"

	"foxhunt.dev/fox/internal"
	"foxhunt.dev/fox/internal/pkg/data/stream"
	"foxhunt.dev/fox/internal/pkg/types/event"
)

const LocalHost = "http://localhost:8080"

type Ecs struct {
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

	// internal
	url string
}

func New(url string) Ecs {
	ecs := Ecs{url: url}
	ecs.Ecs.Version = "9.1.0"
	ecs.Agent.Version = res.Version
	ecs.Agent.Type = "fox"
	ecs.Event.Kind = "event"
	return ecs
}

func (ecs Ecs) String() string {
	return fmt.Sprintf("ECS @ %s", ecs.url)
}

func (ecs Ecs) Stream(e *event.Event) error {
	cef := e.ToCEF()

	// basic properties
	ecs.Timestamp = e.Time.UTC()
	ecs.Message = e.Message
	ecs.Host.Hostname = e.Host
	ecs.User.ID = e.User

	// original event
	ecs.Event.ID = e.Sequence
	ecs.Event.Module = e.Source
	ecs.Event.Dataset = fmt.Sprintf("%s.%s", e.Source, e.Category)
	ecs.Event.Provider = e.Service
	ecs.Event.Severity = int64(e.Severity)
	ecs.Event.Ingested = time.Now().UTC()
	ecs.Event.Original = cef
	ecs.Event.Hash = fmt.Sprintf("%x", murmur3.StringSum64(cef))

	// os specific
	if e.Source == types.Eventlog {
		ecs.Event.Code = e.Fields["EventID"]
	}

	// add fields
	ecs.Labels = make(map[string]any)

	for k, v := range e.Fields {
		ecs.Labels[k] = v
	}

	buf, err := json.Marshal(ecs)

	if err != nil {
		return err
	}

	return stream.Post(ecs.url, string(buf), map[string]string{
		"Content-Type": "application/json",
	})
}
