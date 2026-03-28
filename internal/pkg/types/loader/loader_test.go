package loader

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	_zip "go.foxforensics.dev/fox/v4/internal/pkg/file/archive/7z"

	"go.foxforensics.dev/fox/v4/internal/pkg/file/archive/ar"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/archive/cab"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/archive/cpio"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/archive/iso"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/archive/rar"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/archive/rpm"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/archive/tar"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/archive/xar"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/archive/zip"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/convert/bin/elf"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/convert/bin/ese"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/convert/bin/lnk"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/convert/bin/pe"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/convert/bin/pf"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/convert/log/evtx"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/convert/log/journal"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/deflate/bgzf"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/deflate/br"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/deflate/bzip2"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/deflate/gzip"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/deflate/kanzi"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/deflate/lz4"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/deflate/lzfse"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/deflate/lzip"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/deflate/lznt1"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/deflate/lzo"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/deflate/lzw"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/deflate/minlz"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/deflate/s2"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/deflate/snappy"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/deflate/xz"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/deflate/zlib"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/deflate/zstd"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/format/json"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/format/jsonl"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/format/xml"
	"go.foxforensics.dev/fox/v4/internal/pkg/test"
	"go.foxforensics.dev/fox/v4/internal/pkg/types"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/register"
)

func TestMain(m *testing.M) {
	register.Deflate("bgzf", bgzf.Detect, bgzf.Deflate)
	register.Deflate("br", br.Detect, br.Deflate)
	register.Deflate("bzip2", bzip2.Detect, bzip2.Deflate)
	register.Deflate("gzip", gzip.Detect, gzip.Deflate)
	register.Deflate("kanzi", kanzi.Detect, kanzi.Deflate)
	register.Deflate("lz4", lz4.Detect, lz4.Deflate)
	register.Deflate("lzip", lzip.Detect, lzip.Deflate)
	register.Deflate("lzo", lzo.Detect, lzo.Deflate)
	register.Deflate("lzfse", lzfse.Detect, lzfse.Deflate)
	register.Deflate("lznt1", lznt1.Detect, lznt1.Deflate)
	register.Deflate("lzw", lzw.Detect, lzw.Deflate)
	register.Deflate("minlz", minlz.Detect, minlz.Deflate)
	register.Deflate("s2", s2.Detect, s2.Deflate)
	register.Deflate("snappy", snappy.Detect, snappy.Deflate)
	register.Deflate("xz", xz.Detect, xz.Deflate)
	register.Deflate("zlib", zlib.Detect, zlib.Deflate)
	register.Deflate("zstd", zstd.Detect, zstd.Deflate)

	register.Archive("7z", _zip.Detect, _zip.Extract)
	register.Archive("ar", ar.Detect, ar.Extract)
	register.Archive("cab", cab.Detect, cab.Extract)
	register.Archive("cpio", cpio.Detect, cpio.Extract)
	register.Archive("iso", iso.Detect, iso.Extract)
	register.Archive("rar", rar.Detect, rar.Extract)
	register.Archive("rpm", rpm.Detect, rpm.Extract)
	register.Archive("tar", tar.Detect, tar.Extract)
	register.Archive("xar", xar.Detect, xar.Extract)
	register.Archive("zip", zip.Detect, zip.Extract)

	register.Convert("elf", elf.Detect, elf.Convert)
	register.Convert("ese", ese.Detect, ese.Convert)
	register.Convert("lnk", lnk.Detect, lnk.Convert)
	register.Convert("pe", pe.Detect, pe.Convert)
	register.Convert("pf", pf.Detect, pf.Convert)
	register.Convert("evtx", evtx.Detect, evtx.Convert)
	register.Convert("journal", journal.Detect, journal.Convert)

	register.Format("json", json.Detect, json.Format)
	register.Format("jsonl", jsonl.Detect, jsonl.Format)
	register.Format("xml", xml.Detect, xml.Format)

	os.Exit(m.Run())
}

func TestLoadFiles(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   []string
		out  []string
	}{
		{
			"Empty file",
			[]string{
				"misc/empty",
			},
			[]string{
				"empty",
			},
		}, {
			"Single file",
			[]string{
				"text/bible.txt",
			},
			[]string{
				"bible.txt",
			},
		}, {
			"Multiple files",
			[]string{
				"misc/mbr.bin",
				"misc/rnd.bin",
			},
			[]string{
				"mbr.bin",
				"rnd.bin",
			},
		}, {
			"Multiple file streams",
			[]string{
				"archive/test.zip",
			},
			[]string{
				"test.zip:hello.rar:hello.txt",
				"test.zip:hello.txt.bz2",
				"test.zip:hello.txt.gz",
				"test.zip:hello.txt.lz4",
				"test.zip:hello.txt.xz",
				"test.zip:hello.txt.zst",
			},
		}, {
			"Directory",
			[]string{
				"rules",
			},
			[]string{
				"test.yml",
			},
		}, {
			"Globbing",
			[]string{
				"**/*.txt",
			},
			[]string{
				"bible.txt",
				"fox.txt",
				"nasty.txt",
				"strings.txt",
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			l := New(newOpts())

			paths := consume(l, test.FixtureDir(tt.in))

			if len(paths) != len(tt.out) {
				t.Fatal("invalid count")
			}

			if !slices.Equal(paths, tt.out) {
				t.Fatal("invalid paths")
			}
		})
	}
}

func newOpts() *Options {
	return &Options{
		&types.Limits{},
		&types.Filters{},
		"",
		"",
		1,
		0,
		false,
	}
}

func consume(ldr *Loader, in []string) (out []string) {
	for h := range ldr.Load(in) {
		out = append(out, filepath.Base(h.Name))
		h.Discard()
	}

	slices.Sort(out)

	return out
}
