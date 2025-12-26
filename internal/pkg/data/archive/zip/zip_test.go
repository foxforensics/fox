package zip

import (
	"testing"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

const pass = "test"
const file1 = "archive/fox.zip"
const file2 = "archive/fox.enc.zip"

func BenchmarkDetect(b *testing.B) {
	buf := data.Fixture(file1)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkExtract(b *testing.B) {
	buf := data.Fixture(file1)

	for b.Loop() {
		Extract(buf, "", "")
	}
}

func TestDetect(t *testing.T) {
	if !Detect(data.Fixture(file1)) {
		t.Fatal("not detected")
	}
}

func TestExtract(t *testing.T) {
	for _, tt := range []struct {
		file, pass string
	}{
		{file1, ""},
		{file2, pass},
	} {
		t.Run(tt.file, func(t *testing.T) {
			e := Extract(data.Fixture(tt.file), "", tt.pass)

			if len(e) != 1 {
				t.Fatal("invalid entry count")
			}

			if e[0].Path != data.JoinPart("", data.Sample) {
				t.Fatal("invalid entry path")
			}

			if !data.Assert(e[0].Data) {
				t.Fatal("invalid entry data")
			}
		})
	}
}
