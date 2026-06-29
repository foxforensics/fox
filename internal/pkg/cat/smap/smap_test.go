package smap

import (
	"testing"

	"github.com/dlclark/regexp2/v2"
	"go.foxforensics.eu/fox/v4/internal/test"
)

func BenchmarkMap(b *testing.B) {
	v := test.Fixture("texts/bible.txt")

	for b.Loop() {
		Map(v)
	}
}

func BenchmarkGrep(b *testing.B) {
	v := test.Fixture("texts/bible.txt")

	s := Map(v)

	re := regexp2.MustCompile(".*")

	for b.Loop() {
		s.Grep(re, 2)
	}
}

func TestMap(t *testing.T) {
	v := test.Fixture("texts/bible.txt")

	if len(Map(v)) != 31107 {
		t.Fatal("wrong size")
	}
}

func TestGrep(t *testing.T) {
	v := test.Fixture("texts/bible.txt")

	re := regexp2.MustCompile("King James")

	s := Map(v).Grep(re, 2)

	if len(s) != 1 {
		t.Fatal("wrong length")
	}

	if string(s[0].Bytes) != "Authorized King James Version" {
		t.Fatal("wrong string")
	}
}
