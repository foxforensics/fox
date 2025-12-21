package register

import "github.com/cuhsat/fox/v4/internal/pkg/data"

var Deflates []DeflateEntry
var Archives []ArchiveEntry
var Formats []FormatEntry

type DeflateEntry struct {
	Name    string
	Detect  data.Detect
	Deflate data.Deflate
}

type ArchiveEntry struct {
	Name    string
	Detect  data.Detect
	Extract data.Extract
}

type FormatEntry struct {
	Name   string
	Detect data.Detect
	Format data.Format
}

func Deflate(s string, fn1 data.Detect, fn2 data.Deflate) {
	Deflates = append(Deflates, DeflateEntry{s, fn1, fn2})
}

func Archive(s string, fn1 data.Detect, fn2 data.Extract) {
	Archives = append(Archives, ArchiveEntry{s, fn1, fn2})
}

func Format(s string, fn1 data.Detect, fn2 data.Format) {
	Formats = append(Formats, FormatEntry{s, fn1, fn2})
}
