package formats

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"go.foxforensics.eu/fox/v4/internal/pkg/types/event"
	"go.foxforensics.eu/fox/v4/internal/sys/terminal"
)

const Prefix = ""
const Indent = "  "

func Auto(a fmt.Stringer, json, jsonl bool) string {
	switch {
	case jsonl:
		return terminal.ColorizeAs(AsJSONL(a), "json")
	case json:
		return terminal.ColorizeAs(AsJSON(a), "json")
	default:
		return a.String()
	}
}

func Event(e *event.Event, json, jsonl bool) string {
	if json || jsonl {
		return Auto(e, json, jsonl)
	}

	return terminal.MarkEvent(e.AsCEF())
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
