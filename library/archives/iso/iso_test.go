package iso

import (
	"testing"

	"go.foxforensics.eu/fox/v5/internal/pkg"
	"go.foxforensics.eu/fox/v5/internal/pkg/tests"
)

const src = "archives/fox.iso"

func BenchmarkDetect(b *testing.B) {
	buf := tests.Fixture(src)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkExtract(b *testing.B) {
	buf := tests.Fixture(src)

	for b.Loop() {
		Extract(buf, "", "")
	}
}

func TestDetect(t *testing.T) {
	if !Detect(tests.Fixture(src)) {
		t.Fatal("not detected")
	}
}

func TestExtract(t *testing.T) {
	e := Extract(tests.Fixture(src), "", "")

	if len(e) != 1 {
		t.Fatal("invalid entry count")
	}

	if e[0].Path != pkg.JoinPart("", tests.Fox) {
		t.Fatal("invalid entry path", e[0].Path)
	}

	if !tests.AssertFox(e[0].Data) {
		t.Fatal("invalid entry data")
	}
}
