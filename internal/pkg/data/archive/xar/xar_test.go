package xar

import (
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

const file = "archive/fox.xar"

func BenchmarkDetect(b *testing.B) {
	buf := data.Fixture(file)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkExtract(b *testing.B) {
	buf := data.Fixture(file)

	for b.Loop() {
		Extract(buf, "", "")
	}
}

func TestDetect(t *testing.T) {
	if !Detect(data.Fixture(file)) {
		t.Fatal("not detected")
	}
}

func TestExtract(t *testing.T) {
	e := Extract(data.Fixture(file), "", "")

	if len(e) != 1 {
		t.Fatal("invalid entry count")
	}

	if e[0].Path != data.JoinPart("", data.Sample) {
		t.Fatal("invalid entry path")
	}

	if !data.Assert(e[0].Data) {
		t.Fatal("invalid entry data")
	}
}
