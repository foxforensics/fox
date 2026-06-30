package formats

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"go.foxforensics.eu/fox/v4/internal/pkg/hunt/event"
	"go.foxforensics.eu/fox/v4/internal/sys/writer"
)

const Prefix = ""
const Indent = "  "

func Auto(a fmt.Stringer, json, jsonl bool) string {
	switch {
	case jsonl:
		return writer.ColorizeAs(AsJSONL(a), "json")
	case json:
		return writer.ColorizeAs(AsJSON(a), "json")
	default:
		return a.String()
	}
}

func Event(e *event.Event, json, jsonl bool) string {
	if json || jsonl {
		return Auto(e, json, jsonl)
	}

	return writer.MarkEvent(e.AsCEF())
}

func AsJSON(a any) string {
	b, err := json.MarshalIndent(a, Prefix, Indent)

	if err != nil {
		slog.Error(err.Error())
		return JsonError(err)
	}

	return string(b)
}

func AsJSONL(a any) string {
	b, err := json.Marshal(a)

	if err != nil {
		slog.Error(err.Error())
		return JsonError(err)
	}

	return string(b)
}

func JsonError(err error) string {
	return fmt.Sprintf(`{"error": "%s"}`, err.Error())
}
