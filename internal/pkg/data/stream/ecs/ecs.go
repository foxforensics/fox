// Package ecs specification:
// https://www.elastic.co/docs/reference/ecs/ecs-field-reference
package ecs

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cuhsat/fox/v4/internal"
	"github.com/cuhsat/fox/v4/internal/pkg/data/stream"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
)

const version = "9.1.0"

type Ecs struct {
	stream.Stream

	Timestamp time.Time `json:"@timestamp"`
	Message   string    `json:"message"`

	Agent struct {
		Type    string `json:"type"`
		Version string `json:"version"`
	} `json:"agent"`

	Ecs struct {
		Version string `json:"version"`
	} `json:"ecs"`
}

type Evt struct {
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
}

func New(url string) Ecs {
	ecs := Ecs{Stream: stream.Stream{Url: url, Map: map[string]string{
		"Content-Type": "application/json",
	}}}

	ecs.Ecs.Version = version
	ecs.Agent.Type = "fox"
	ecs.Agent.Version = app.Version[1:]

	return ecs
}

func (ecs Ecs) String() string {
	return fmt.Sprintf("ECS: %s", ecs.Url)
}

func (ecs Ecs) Write(e *event.Event) (int64, int64, error) {
	ecs.Timestamp = time.Now().UTC()
	ecs.Message = strings.TrimRight(e.ToCEF(), "\n")

	buf, err := json.Marshal(ecs)

	if err != nil {
		return 0, 0, nil
	}

	return ecs.Post(string(buf))
}
