package evtx

import (
	"bufio"
	"bytes"
	"encoding/json"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

const file = "convert/test.evtx"

func BenchmarkDetect(b *testing.B) {
	buf := data.Fixture(file)

	b.ResetTimer()

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkConvert(b *testing.B) {
	buf := data.Fixture(file)

	b.ResetTimer()

	for b.Loop() {
		_, _ = Convert(buf)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(data.Fixture(file)) {
		t.Fatal("not detected")
	}
}

func TestConvert(t *testing.T) {
	buf, err := Convert(data.Fixture(file))

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
