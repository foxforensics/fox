package register

import (
	"sync"

	"go.foxforensics.eu/fox/v4/internal/pkg"
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
	Detect pkg.Detect
	Format pkg.Format
}

type DeflateEntry struct {
	Name    string
	Detect  pkg.Detect
	Deflate pkg.Deflate
}

type ExtractEntry struct {
	Name    string
	Detect  pkg.Detect
	Extract pkg.Extract
}

type ConvertEntry struct {
	Name    string
	Detect  pkg.Detect
	Convert pkg.Convert
}

func Format(s string, fn1 pkg.Detect, fn2 pkg.Format) {
	Registry.Lock()
	Registry.Formats = append(Registry.Formats, FormatEntry{s, fn1, fn2})
	Registry.Unlock()
}

func Deflate(s string, fn1 pkg.Detect, fn2 pkg.Deflate) {
	Registry.Lock()
	Registry.Deflates = append(Registry.Deflates, DeflateEntry{s, fn1, fn2})
	Registry.Unlock()
}

func Extract(s string, fn1 pkg.Detect, fn2 pkg.Extract) {
	Registry.Lock()
	Registry.Extracts = append(Registry.Extracts, ExtractEntry{s, fn1, fn2})
	Registry.Unlock()
}

func Convert(s string, fn1 pkg.Detect, fn2 pkg.Convert) {
	Registry.Lock()
	Registry.Converts = append(Registry.Converts, ConvertEntry{s, fn1, fn2})
	Registry.Unlock()
}
