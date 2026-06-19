package loader

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"go.foxforensics.eu/fox/v4/internal/pkg/types"
	"go.foxforensics.eu/fox/v4/internal/test"
)

func TestMain(m *testing.M) {
	RegisterDeflates()
	RegisterExtracts()
	RegisterConverts()
	RegisterFormats()

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
				"binary/test.nil",
			},
			[]string{
				"test.nil",
			},
		}, {
			"Single file",
			[]string{
				"string/bible.txt",
			},
			[]string{
				"bible.txt",
			},
		}, {
			"Multiple files",
			[]string{
				"binary/test.mbr",
				"binary/test.rnd",
			},
			[]string{
				"test.mbr",
				"test.rnd",
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
				"string",
			},
			[]string{
				"bible.txt",
				"nasty.txt",
				"test.txt",
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
				"test.txt",
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
		1,
		false,
	}
}

func consume(ldr *Loader, in []string) (out []string) {
	for h := range ldr.Load(context.Background(), in) {
		out = append(out, filepath.Base(h.Name))
		h.Discard()
	}

	slices.Sort(out)

	return out
}
