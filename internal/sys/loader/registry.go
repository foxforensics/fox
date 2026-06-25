package loader

import (
	"go.foxforensics.eu/fox/v4/internal/pkg"
	_zip "go.foxforensics.eu/fox/v4/internal/pkg/adapters/archives/7z"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/archives/ar"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/archives/cab"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/archives/cpio"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/archives/iso"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/archives/msi"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/archives/rar"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/archives/rpm"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/archives/tar"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/archives/xar"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/archives/zip"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/binaries/bin/elf"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/binaries/bin/ese"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/binaries/bin/lnk"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/binaries/bin/mft"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/binaries/bin/pe"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/binaries/bin/pf"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/binaries/bin/pst"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/binaries/log/evtx"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/binaries/log/fortinet"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/binaries/log/journal"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/deflates/bgzf"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/deflates/br"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/deflates/bzip2"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/deflates/gzip"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/deflates/kanzi"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/deflates/lz4"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/deflates/lzfse"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/deflates/lzip"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/deflates/lznt1"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/deflates/lzo"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/deflates/lzw"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/deflates/minlz"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/deflates/s2"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/deflates/snappy"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/deflates/xz"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/deflates/zlib"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/deflates/zstd"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/formats/json"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/formats/jsonl"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/formats/xml"
)

var registry = &struct {
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
}

func RegisterExtracts() {
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
}

func RegisterConverts() {
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
}

func RegisterFormats() {
	registry.Formats = []FormatEntry{
		{"json", json.Detect, json.Format},
		{"jsonl", jsonl.Detect, jsonl.Format},
		{"xml", xml.Detect, xml.Format},
	}
}
