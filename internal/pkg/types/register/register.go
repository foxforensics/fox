package register

import "github.com/cuhsat/fox/v4/internal/pkg/data"

var (
	Readers  []ReaderEntry
	Formats  []FormatEntry
	Deflates []DeflateEntry
	Archives []ArchiveEntry
	Converts []ConvertEntry
)

type ReaderEntry struct {
	Name   string
	Detect data.Detect
	Reader data.Reader
}

type FormatEntry struct {
	Name   string
	Detect data.Detect
	Format data.Format
}

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

type ConvertEntry struct {
	Name    string
	Detect  data.Detect
	Convert data.Convert
}

func Reader(s string, fn1 data.Detect, fn2 data.Reader) {
	Readers = append(Readers, ReaderEntry{s, fn1, fn2})
}

func Format(s string, fn1 data.Detect, fn2 data.Format) {
	Formats = append(Formats, FormatEntry{s, fn1, fn2})
}

func Deflate(s string, fn1 data.Detect, fn2 data.Deflate) {
	Deflates = append(Deflates, DeflateEntry{s, fn1, fn2})
}

func Archive(s string, fn1 data.Detect, fn2 data.Extract) {
	Archives = append(Archives, ArchiveEntry{s, fn1, fn2})
}

func Convert(s string, fn1 data.Detect, fn2 data.Convert) {
	Converts = append(Converts, ConvertEntry{s, fn1, fn2})
}
