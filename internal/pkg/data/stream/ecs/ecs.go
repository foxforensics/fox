// Package ecs specification:
// https://www.elastic.co/docs/reference/ecs/ecs-field-reference
package ecs

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/cuhsat/fox/v4/internal"
	"github.com/cuhsat/fox/v4/internal/pkg/data/stream"
)

const version = "9.1.0"

type Ecs struct {
	stream.Schema

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

func New(url string) Ecs {
	ecs := Ecs{Schema: stream.Schema{Url: url, Map: map[string]string{
		"Content-Type": "application/json",
	}}}

	ecs.Ecs.Version = version
	ecs.Agent.Type = "fox"
	ecs.Agent.Version = app.Version[1:]

	return ecs
}

func (ecs Ecs) Write(p []byte) (int, error) {
	ecs.Timestamp = time.Now().UTC()
	ecs.Message = strings.TrimRight(string(p), "\n")

	buf, err := json.Marshal(ecs)

	if err != nil {
		return 0, nil
	}

	return ecs.Post(string(buf))
}
