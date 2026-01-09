package ewf

import (
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

const file = "image/test.E01"

func BenchmarkDetect(b *testing.B) {
	buf := data.Fixture(file)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkIngest(b *testing.B) {
	buf := data.Fixture(file)

	for b.Loop() {
		_, _ = Ingest(buf)
	}
}

func TestDetect(t *testing.T) {
	if !Detect(data.Fixture(file)) {
		t.Fatal("not detected")
	}
}

func TestIngest(t *testing.T) {
	buf, err := Ingest(data.Fixture(file))

	if err != nil {
		t.Error(err)
	}

	if len(buf) == 0 {
		t.Fatal("not ingested")
	}
}
