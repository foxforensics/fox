package ar

import (
	"testing"

	"foxhunt.dev/fox/internal/pkg/data"
	"foxhunt.dev/fox/internal/pkg/test"
)

const file = "archive/fox.ar"

func BenchmarkDetect(b *testing.B) {
	buf := test.Fixture(file)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkExtract(b *testing.B) {
	buf := test.Fixture(file)

	for b.Loop() {
		Extract(buf, "", "")
	}
}

func TestDetect(t *testing.T) {
	if !Detect(test.Fixture(file)) {
		t.Fatal("not detected")
	}
}

func TestExtract(t *testing.T) {
	e := Extract(test.Fixture(file), "", "")

	if len(e) != 1 {
		t.Fatal("invalid entry count")
	}

	if e[0].Path != data.JoinPart("", test.Sample) {
		t.Fatal("invalid entry path")
	}

	if !test.Assert(e[0].Data) {
		t.Fatal("invalid entry data")
	}
}
