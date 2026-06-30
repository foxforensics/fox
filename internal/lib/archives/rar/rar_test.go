package rar

import (
	"testing"

	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/fox/v4/internal/test"
)

const pass = "test"
const src1 = "archives/fox.rar"
const src2 = "archives/fox.crypt.rar"

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

			if e[0].Path != sys.JoinPart("", test.Fox) {
				t.Fatal("invalid entry path")
			}

			if !test.AssertFox(e[0].Data) {
				t.Fatal("invalid entry data")
			}
		})
	}
}
