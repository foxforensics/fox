package loader

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	_zip "foxhunt.dev/fox/internal/pkg/data/archive/7z"

	"foxhunt.dev/fox/internal/pkg/data/archive/ar"
	"foxhunt.dev/fox/internal/pkg/data/archive/cab"
	"foxhunt.dev/fox/internal/pkg/data/archive/cpio"
	"foxhunt.dev/fox/internal/pkg/data/archive/iso"
	"foxhunt.dev/fox/internal/pkg/data/archive/rar"
	"foxhunt.dev/fox/internal/pkg/data/archive/rpm"
	"foxhunt.dev/fox/internal/pkg/data/archive/tar"
	"foxhunt.dev/fox/internal/pkg/data/archive/xar"
	"foxhunt.dev/fox/internal/pkg/data/archive/zip"
	"foxhunt.dev/fox/internal/pkg/data/convert/bin/elf"
	"foxhunt.dev/fox/internal/pkg/data/convert/bin/ese"
	"foxhunt.dev/fox/internal/pkg/data/convert/bin/lnk"
	"foxhunt.dev/fox/internal/pkg/data/convert/bin/pe"
	"foxhunt.dev/fox/internal/pkg/data/convert/bin/pf"
	"foxhunt.dev/fox/internal/pkg/data/convert/log/evtx"
	"foxhunt.dev/fox/internal/pkg/data/convert/log/fortinet"
	"foxhunt.dev/fox/internal/pkg/data/convert/log/journal"
	"foxhunt.dev/fox/internal/pkg/data/deflate/br"
	"foxhunt.dev/fox/internal/pkg/data/deflate/bzip2"
	"foxhunt.dev/fox/internal/pkg/data/deflate/gzip"
	"foxhunt.dev/fox/internal/pkg/data/deflate/kanzi"
	"foxhunt.dev/fox/internal/pkg/data/deflate/lz4"
	"foxhunt.dev/fox/internal/pkg/data/deflate/lzfse"
	"foxhunt.dev/fox/internal/pkg/data/deflate/lzip"
	"foxhunt.dev/fox/internal/pkg/data/deflate/lzo"
	"foxhunt.dev/fox/internal/pkg/data/deflate/lzw"
	"foxhunt.dev/fox/internal/pkg/data/deflate/minlz"
	"foxhunt.dev/fox/internal/pkg/data/deflate/s2"
	"foxhunt.dev/fox/internal/pkg/data/deflate/snappy"
	"foxhunt.dev/fox/internal/pkg/data/deflate/xz"
	"foxhunt.dev/fox/internal/pkg/data/deflate/zlib"
	"foxhunt.dev/fox/internal/pkg/data/deflate/zstd"
	"foxhunt.dev/fox/internal/pkg/data/format/json"
	"foxhunt.dev/fox/internal/pkg/data/format/jsonl"
	"foxhunt.dev/fox/internal/pkg/test"
	"foxhunt.dev/fox/internal/pkg/types"
	"foxhunt.dev/fox/internal/pkg/types/register"
)

func TestMain(m *testing.M) {
	register.Deflate("br", br.Detect, br.Deflate)
	register.Deflate("bzip2", bzip2.Detect, bzip2.Deflate)
	register.Deflate("gzip", gzip.Detect, gzip.Deflate)
	register.Deflate("kanzi", kanzi.Detect, kanzi.Deflate)
	register.Deflate("lz4", lz4.Detect, lz4.Deflate)
	register.Deflate("lzip", lzip.Detect, lzip.Deflate)
	register.Deflate("lzo", lzo.Detect, lzo.Deflate)
	register.Deflate("lzfse", lzfse.Detect, lzfse.Deflate)
	register.Deflate("lzw", lzw.Detect, lzw.Deflate)
	register.Deflate("minlz", minlz.Detect, minlz.Deflate)
	register.Deflate("s2", s2.Detect, s2.Deflate)
	register.Deflate("snappy", snappy.Detect, snappy.Deflate)
	register.Deflate("xz", xz.Detect, xz.Deflate)
	register.Deflate("zlib", zlib.Detect, zlib.Deflate)
	register.Deflate("zstd", zstd.Detect, zstd.Deflate)

	register.Archive("ar", ar.Detect, ar.Extract)
	register.Archive("cab", cab.Detect, cab.Extract)
	register.Archive("cpio", cpio.Detect, cpio.Extract)
	register.Archive("iso", iso.Detect, iso.Extract)
	register.Archive("rar", rar.Detect, rar.Extract)
	register.Archive("rpm", rpm.Detect, rpm.Extract)
	register.Archive("7z", _zip.Detect, _zip.Extract)
	register.Archive("tar", tar.Detect, tar.Extract)
	register.Archive("xar", xar.Detect, xar.Extract)
	register.Archive("zip", zip.Detect, zip.Extract)

	register.Convert("elf", elf.Detect, elf.Convert)
	register.Convert("ese", ese.Detect, ese.Convert)
	register.Convert("lnk", lnk.Detect, lnk.Convert)
	register.Convert("pe", pe.Detect, pe.Convert)
	register.Convert("pf", pf.Detect, pf.Convert)
	register.Convert("evtx", evtx.Detect, evtx.Convert)
	register.Convert("fortinet", fortinet.Detect, fortinet.Convert)
	register.Convert("journal", journal.Detect, journal.Convert)

	register.Format("json", json.Detect, json.Format)
	register.Format("jsonl", jsonl.Detect, jsonl.Format)

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
				"archive/multi.zip",
			},
			[]string{
				"multi.zip:hello.rar:hello.txt",
				"multi.zip:hello.txt.bz2",
				"multi.zip:hello.txt.gz",
				"multi.zip:hello.txt.lz4",
				"multi.zip:hello.txt.xz",
				"multi.zip:hello.txt.zst",
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
				"dump.txt",
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
