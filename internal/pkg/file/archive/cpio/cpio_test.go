package cpio

import (
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/file"
	"github.com/cuhsat/fox/v4/internal/pkg/test"
)

const archive = "archive/fox.cpio"

func BenchmarkDetect(b *testing.B) {
	buf := test.Fixture(archive)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkExtract(b *testing.B) {
	buf := test.Fixture(archive)

	for b.Loop() {
		Extract(buf, "", "")
	}
}

func TestDetect(t *testing.T) {
	if !Detect(test.Fixture(archive)) {
		t.Fatal("not detected")
	}
}

func TestExtract(t *testing.T) {
	e := Extract(test.Fixture(archive), "", "")

	if len(e) != 1 {
		t.Fatal("invalid entry count")
	}

	if e[0].Path != file.JoinPart("", test.Sample) {
		t.Fatal("invalid entry path")
	}

	if !test.Assert(e[0].Data) {
		t.Fatal("invalid entry data")
	}
}
