package formats

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"go.foxforensics.eu/fox/v5/internal/pkg/hunt/event"
	"go.foxforensics.eu/fox/v5/internal/sys/writer"
)

const Prefix = ""
const Indent = "  "

func Auto(s fmt.Stringer, json, jsonl bool) string {
	switch {
	case json:
		return writer.ColorizeAs(AsJSON(s), "json")
	case jsonl:
		return writer.ColorizeAs(AsJSONL(s), "json")
	default:
		return s.String()
	}
}

func Event(e *event.Event, json, jsonl bool) string {
	switch {
	case json:
		return writer.ColorizeAs(AsJSON(e), "json")
	case jsonl:
		return writer.ColorizeAs(AsJSONL(e), "json")
	default:
		return writer.MarkEvent(e.AsCEF())
	}
}

func AsJSON(a any) string {
	b, err := json.MarshalIndent(a, Prefix, Indent)

	if err != nil {
		slog.Error(err.Error())
		return string(JsonError(err))
	}

	return string(b)
}

func AsJSONL(a any) string {
	b, err := json.Marshal(a)

	if err != nil {
		slog.Error(err.Error())
		return string(JsonError(err))
	}

	return string(b)
}

func JsonError(e error) []byte {
	b, err := json.Marshal(struct {
		Error string `json:"error"`
	}{
		e.Error(), // escape error
	})

	if err != nil {
		slog.Error(err.Error())
		return []byte(`{"error": "unknown error"}`)
	}

	return b
}
