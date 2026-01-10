package loader

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"

	szip "github.com/cuhsat/fox/v4/internal/pkg/data/archive/7z"

	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/ar"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/cab"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/cpio"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/rar"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/rpm"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/tar"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/xar"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/zip"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/bin/pe"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/log/evtx"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/log/journal"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/br"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/bzip2"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/gzip"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/kanzi"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/lz4"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/lzip"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/lzw"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/minlz"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/s2"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/snappy"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/xz"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/zlib"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/zstd"
	"github.com/cuhsat/fox/v4/internal/pkg/data/image/ewf"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/register"
)

func TestMain(m *testing.M) {
	register.Deflate("br", br.Detect, br.Deflate)
	register.Deflate("bzip2", bzip2.Detect, bzip2.Deflate)
	register.Deflate("gzip", gzip.Detect, gzip.Deflate)
	register.Deflate("kanzi", kanzi.Detect, kanzi.Deflate)
	register.Deflate("lz4", lz4.Detect, lz4.Deflate)
	register.Deflate("lzip", lzip.Detect, lzip.Deflate)
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
	register.Archive("rar", rar.Detect, rar.Extract)
	register.Archive("rpm", rpm.Detect, rpm.Extract)
	register.Archive("szip", szip.Detect, szip.Extract)
	register.Archive("tar", tar.Detect, tar.Extract)
	register.Archive("xar", xar.Detect, xar.Extract)
	register.Archive("zip", zip.Detect, zip.Extract)

	register.Convert("evtx", evtx.Detect, evtx.Convert)
	register.Convert("journal", journal.Detect, journal.Convert)
	register.Convert("pe", pe.Detect, pe.Convert)

	register.Image("ewf", ewf.Detect, ewf.Ingest)

	os.Exit(m.Run())
}

func TestLoadInput(t *testing.T) {
	opts := newOpts()
	opts.Input = "Test"

	l := New(opts)

	paths := consume(l, []string{})

	if len(paths) != 1 {
		t.Fatal("invalid count")
	}

	if paths[0] != "input" {
		t.Fatal("invalid path")
	}
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
				"misc",
			},
			[]string{
				"eicar",
				"empty",
				"mbr.bin",
				"rnd.bin",
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

			paths := consume(l, fixture(tt.in))

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
		"",
		1,
		0,
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

func fixture(in []string) (out []string) {
	const dir = "../../../../testdata"

	_, c, _, ok := runtime.Caller(0)

	if !ok {
		log.Fatalln("runtime error")
	}

	for _, path := range in {
		out = append(out, filepath.Join(filepath.Dir(c), dir, path))
	}

	return
}
