package format

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

const Prefix = ""
const Indent = "  "

func Error(err error) string {
	return fmt.Sprintf(`{"error": "%s"}`, err.Error())
}

func AsJSON(a any) string {
	b, err := json.MarshalIndent(a, Prefix, Indent)

	if err != nil {
		slog.Error(err.Error())
		return Error(err)
	}

	return string(b)
}

func AsJSONL(a any) string {
	b, err := json.Marshal(a)

	if err != nil {
		slog.Error(err.Error())
		return Error(err)
	}

	return string(b)
}
