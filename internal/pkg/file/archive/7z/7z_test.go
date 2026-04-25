package szip

import (
	"testing"

	"go.foxforensics.dev/fox/v4/internal/pkg/file"
	"go.foxforensics.dev/fox/v4/internal/pkg/test"
)

const pass = "test"
const src1 = "archive/fox.7z"
const src2 = "archive/crypt.7z"

func BenchmarkDetect(b *testing.B) {
	buf := test.Fixture(src1)

	for b.Loop() {
		_ = Detect(buf)
	}
}

func BenchmarkExtract(b *testing.B) {
	buf := test.Fixture(src1)

	for b.Loop() {
		Extract(buf, "", "")
	}
}

func TestDetect(t *testing.T) {
	if !Detect(test.Fixture(src1)) {
		t.Fatal("not detected")
	}
}

func TestExtract(t *testing.T) {
	for _, tt := range []struct {
		file, pass string
	}{
		{src1, ""},
		{src2, pass},
	} {
		t.Run(tt.file, func(t *testing.T) {
			e := Extract(test.Fixture(tt.file), "", tt.pass)

			if len(e) != 1 {
				t.Fatal("invalid entry count")
			}

			if e[0].Path != file.JoinPart("", test.Sample) {
				t.Fatal("invalid entry path")
			}

			if !test.Assert(e[0].Data) {
				t.Fatal("invalid entry data")
			}
		})
	}
}
