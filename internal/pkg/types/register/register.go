package register

import "go.foxforensics.dev/fox/v4/internal/pkg/file"

var (
	Formats  []FormatEntry
	Deflates []DeflateEntry
	Extracts []ExtractEntry
	Converts []ConvertEntry
)

type FormatEntry struct {
	Name   string
	Detect file.Detect
	Format file.Format
}

type DeflateEntry struct {
	Name    string
	Detect  file.Detect
	Deflate file.Deflate
}

type ExtractEntry struct {
	Name    string
	Detect  file.Detect
	Extract file.Extract
}

type ConvertEntry struct {
	Name    string
	Detect  file.Detect
	Convert file.Convert
}

func Format(s string, fn1 file.Detect, fn2 file.Format) {
	Formats = append(Formats, FormatEntry{s, fn1, fn2})
}

func Deflate(s string, fn1 file.Detect, fn2 file.Deflate) {
	Deflates = append(Deflates, DeflateEntry{s, fn1, fn2})
}

func Extract(s string, fn1 file.Detect, fn2 file.Extract) {
	Extracts = append(Extracts, ExtractEntry{s, fn1, fn2})
}

func Convert(s string, fn1 file.Detect, fn2 file.Convert) {
	Converts = append(Converts, ConvertEntry{s, fn1, fn2})
}
