package register

import (
	"sync"

	"go.foxforensics.eu/fox/v4/internal/pkg/file"
)

var Registry = &struct {
	sync.Mutex
	Formats  []FormatEntry
	Deflates []DeflateEntry
	Extracts []ExtractEntry
	Converts []ConvertEntry
}{
	Formats:  make([]FormatEntry, 0),
	Deflates: make([]DeflateEntry, 0),
	Extracts: make([]ExtractEntry, 0),
	Converts: make([]ConvertEntry, 0),
}

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
	Registry.Lock()
	Registry.Formats = append(Registry.Formats, FormatEntry{s, fn1, fn2})
	Registry.Unlock()
}

func Deflate(s string, fn1 file.Detect, fn2 file.Deflate) {
	Registry.Lock()
	Registry.Deflates = append(Registry.Deflates, DeflateEntry{s, fn1, fn2})
	Registry.Unlock()
}

func Extract(s string, fn1 file.Detect, fn2 file.Extract) {
	Registry.Lock()
	Registry.Extracts = append(Registry.Extracts, ExtractEntry{s, fn1, fn2})
	Registry.Unlock()
}

func Convert(s string, fn1 file.Detect, fn2 file.Convert) {
	Registry.Lock()
	Registry.Converts = append(Registry.Converts, ConvertEntry{s, fn1, fn2})
	Registry.Unlock()
}
