package journal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

const file = "format/log.journal"

func BenchmarkDetect(b *testing.B) {
	buf := data.Fixture(file)

	b.ResetTimer()

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkFormat(b *testing.B) {
	buf := data.Fixture(file)

	b.ResetTimer()

	for b.Loop() {
		_, _ = Format(buf, 0)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(data.Fixture(file)) {
		t.Fatal("not detected")
	}
}

func TestFormat(t *testing.T) {
	buf, err := Format(data.Fixture(file), 0)

	if err != nil {
		t.Error(err)
	}

	jsonl := bufio.NewScanner(bytes.NewReader(buf))

	for jsonl.Scan() {
		if !json.Valid([]byte(jsonl.Text())) {
			t.Fatal("invalid json")
		}
	}

	if err := jsonl.Err(); err != nil {
		t.Error(err)
	}
}
