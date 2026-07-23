// Package ecs applies this schema: https://www.elastic.co/docs/reference/ecs/ecs-field-reference
package ecs

import (
	"encoding/json"
	"fmt"
	"time"

	"go.foxforensics.eu/fox/v5/internal/cmd/hunt/event"
	"go.foxforensics.eu/fox/v5/internal/pkg/version"
	"go.foxforensics.eu/fox/v5/library/binaries"
	"go.foxforensics.eu/fox/v5/library/formats"
	"go.foxforensics.eu/hasher/hash"
)

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
		ID       string `json:"id,omitempty"`
		Kind     string `json:"kind,omitempty"`
		Module   string `json:"module,omitempty"`
		Dataset  string `json:"dataset,omitempty"`
		Severity int64  `json:"severity,omitempty"`
		Provider string `json:"provider,omitempty"`
		Original string `json:"original,omitempty"`
		Hash     string `json:"hash,omitempty"`
		Code     string `json:"code,omitempty"`
	} `json:"event"`
	Labels    map[string]any `json:"labels,omitempty"`
	Timestamp time.Time      `json:"@timestamp"`
	Message   string         `json:"message"`
}

func Apply(evt *event.Event) ([]byte, error) {
	original := formats.AsJSONL(evt)

	ecs := &Ecs{
		Labels:    make(map[string]any),
		Timestamp: evt.Time.UTC(),
		Message:   evt.Message,
	}

	ecs.Ecs.Version = "9.1.0"
	ecs.Agent.Type = "fox"
	ecs.Agent.Version = version.Number
	ecs.Host.Hostname = evt.Host
	ecs.User.ID = evt.User
	ecs.Event.ID = evt.Sequence
	ecs.Event.Kind = "event"
	ecs.Event.Module = evt.Source
	ecs.Event.Dataset = fmt.Sprintf("%s.%s", evt.Source, evt.Category)
	ecs.Event.Severity = int64(evt.Severity)
	ecs.Event.Provider = evt.Service
	ecs.Event.Original = original
	ecs.Event.Hash = hash.MustSum(hash.MURMUR3, []byte(original))

	if evt.Source == string(binaries.Eventlog) {
		ecs.Event.Code = evt.Fields["EventID"]
	}

	for k, v := range evt.Fields {
		ecs.Labels[k] = v
	}

	return json.Marshal(ecs)
}
