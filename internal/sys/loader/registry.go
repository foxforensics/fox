package loader

import (
	"sync"

	"go.foxforensics.eu/fox/v4/internal/pkg"
	_zip "go.foxforensics.eu/fox/v4/internal/pkg/files/archive/7z"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/archive/ar"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/archive/cab"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/archive/cpio"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/archive/iso"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/archive/msi"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/archive/rar"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/archive/rpm"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/archive/tar"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/archive/xar"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/archive/zip"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/binary/bin/elf"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/binary/bin/ese"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/binary/bin/lnk"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/binary/bin/mft"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/binary/bin/pe"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/binary/bin/pf"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/binary/bin/pst"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/binary/log/evtx"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/binary/log/fortinet"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/binary/log/journal"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/deflate/bgzf"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/deflate/br"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/deflate/bzip2"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/deflate/gzip"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/deflate/kanzi"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/deflate/lz4"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/deflate/lzfse"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/deflate/lzip"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/deflate/lznt1"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/deflate/lzo"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/deflate/lzw"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/deflate/minlz"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/deflate/s2"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/deflate/snappy"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/deflate/xz"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/deflate/zlib"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/deflate/zstd"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/format/json"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/format/jsonl"
	"go.foxforensics.eu/fox/v4/internal/pkg/files/format/xml"
)

var registry = &struct {
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

func RegisterDeflates() {
	registry.Lock()
	registry.Deflates = []DeflateEntry{
		{"bgzf", bgzf.Detect, bgzf.Deflate},
		{"br", br.Detect, br.Deflate},
		{"bzip2", bzip2.Detect, bzip2.Deflate},
		{"gzip", gzip.Detect, gzip.Deflate},
		{"kanzi", kanzi.Detect, kanzi.Deflate},
		{"lz4", lz4.Detect, lz4.Deflate},
		{"lzip", lzip.Detect, lzip.Deflate},
		{"lzo", lzo.Detect, lzo.Deflate},
		{"lzfse", lzfse.Detect, lzfse.Deflate},
		{"lznt1", lznt1.Detect, lznt1.Deflate},
		{"lzw", lzw.Detect, lzw.Deflate},
		{"minlz", minlz.Detect, minlz.Deflate},
		{"s2", s2.Detect, s2.Deflate},
		{"snappy", snappy.Detect, snappy.Deflate},
		{"xz", xz.Detect, xz.Deflate},
		{"zlib", zlib.Detect, zlib.Deflate},
		{"zstd", zstd.Detect, zstd.Deflate},
	}
	registry.Unlock()
}

func RegisterExtracts() {
	registry.Lock()
	registry.Extracts = []ExtractEntry{
		{"7z", _zip.Detect, _zip.Extract},
		{"ar", ar.Detect, ar.Extract},
		{"cab", cab.Detect, cab.Extract},
		{"cpio", cpio.Detect, cpio.Extract},
		{"iso", iso.Detect, iso.Extract},
		{"msi", msi.Detect, msi.Extract},
		{"rar", rar.Detect, rar.Extract},
		{"rpm", rpm.Detect, rpm.Extract},
		{"tar", tar.Detect, tar.Extract},
		{"xar", xar.Detect, xar.Extract},
		{"zip", zip.Detect, zip.Extract},
	}
	registry.Unlock()
}

func RegisterConverts() {
	registry.Lock()
	registry.Converts = []ConvertEntry{
		{"elf", elf.Detect, elf.Convert},
		{"ese", ese.Detect, ese.Convert},
		{"lnk", lnk.Detect, lnk.Convert},
		{"mft", mft.Detect, mft.Convert},
		{"pe", pe.Detect, pe.Convert},
		{"pf", pf.Detect, pf.Convert},
		{"pst", pst.Detect, pst.Convert},
		{"evtx", evtx.Detect, evtx.Convert},
		{"fortinet", fortinet.Detect, fortinet.Convert},
		{"journal", journal.Detect, journal.Convert},
	}
	registry.Unlock()
}

func RegisterFormats() {
	registry.Lock()
	registry.Formats = []FormatEntry{
		{"json", json.Detect, json.Format},
		{"jsonl", jsonl.Detect, jsonl.Format},
		{"xml", xml.Detect, xml.Format},
	}
	registry.Unlock()
}
