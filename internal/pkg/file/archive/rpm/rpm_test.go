package rpm

import (
	"strings"
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/test"
)

const src = "archive/fox.rpm"

func BenchmarkDetect(b *testing.B) {
	buf := test.Fixture(src)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkExtract(b *testing.B) {
	buf := test.Fixture(src)

	for b.Loop() {
		Extract(buf, "", "")
	}
}

func TestDetect(t *testing.T) {
	if !Detect(test.Fixture(src)) {
		t.Fatal("not detected")
	}
}

func TestExtract(t *testing.T) {
	e := Extract(test.Fixture(src), "", "")

	if len(e) != 1 {
		t.Fatal("invalid entry count")
	}

	if !strings.HasSuffix(e[0].Path, test.Sample) {
		t.Fatal("invalid entry path")
	}

	if !test.Assert(e[0].Data) {
		t.Fatal("invalid entry data")
	}
}
