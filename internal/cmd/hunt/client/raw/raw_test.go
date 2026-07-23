package raw

import (
	"testing"

	"go.foxforensics.eu/fox/v5/internal/cmd/hunt/event"
)

func BenchmarkApply(b *testing.B) {
	for b.Loop() {
		_, _ = Apply(&event.Event{})
	}
}

func TestApply(t *testing.T) {
	_, err := Apply(&event.Event{})

	if err != nil {
		t.Error(err)
	}
}
