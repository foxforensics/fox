package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/fatih/color"

	szip "github.com/cuhsat/fox/v4/internal/pkg/data/archive/7z"

	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/ar"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/cab"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/cpio"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/rar"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/rpm"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/tar"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/xar"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/zip"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/evtx"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/journal"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/pe"
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
	"github.com/cuhsat/fox/v4/internal/pkg/data/format/fox"
	"github.com/cuhsat/fox/v4/internal/pkg/data/format/json"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heapset"
	"github.com/cuhsat/fox/v4/internal/pkg/types/register"
	"github.com/cuhsat/fox/v4/internal/pkg/types/writer"
)

type Globals struct {
	// file limits
	Head  bool `short:"h" xor:"head,tail"`
	Tail  bool `short:"t" xor:"head,tail"`
	Lines uint `short:"n" xor:"lines,bytes"`
	Bytes uint `short:"c" xor:"lines,bytes"`

	// file loader
	Input string `short:"i"`
	Pass  string `short:"p"`

	// file writer
	File string `short:"f"`

	// line filter
	Regex   string `short:"e"`
	Context uint   `short:"C"`
	Before  uint   `short:"B"`
	After   uint   `short:"A"`

	// disable
	Raw       bool `short:"r"`
	Quiet     bool `short:"q"`
	NoFile    bool `long:"no-file"`
	NoLine    bool `long:"no-line"`
	NoColor   bool `long:"no-color"`
	NoDeflate bool `long:"no-deflate"`
	NoExtract bool `long:"no-extract"`
	NoConvert bool `long:"no-convert"`

	// standard
	Help    bool
	DryRun  bool `short:"d" long:"dry-run"`
	Verbose int  `short:"v" type:"counter"`

	// bootstrapped
	Stdout io.WriteCloser   `kong:"-"`
	Filter *regexp.Regexp   `kong:"-"`
	Heaps  *heapset.HeapSet `kong:"-"`
}

func (cli *Globals) Load(args []string) *heapset.HeapSet {
	if len(cli.Regex) > 0 {
		cli.Filter = regexp.MustCompile(cli.Regex)
	}

	if cli.Context > 0 {
		cli.Before = cli.Context
		cli.After = cli.Context
	}

	if cli.Raw {
		cli.NoFile = true
		cli.NoLine = true
		cli.NoColor = true
		cli.NoDeflate = true
		cli.NoExtract = true
		cli.NoConvert = true
	}

	if len(cli.File) > 0 {
		cli.NoColor = true
		cli.Stdout = writer.New(cli.File)
	} else if cli.Quiet {
		log.SetOutput(io.Discard)
		cli.Stdout, _ = os.Open(os.DevNull)
	} else {
		cli.Stdout = os.Stdout
	}

	if cli.NoColor {
		color.NoColor = true // turn off color package
	}

	if !cli.NoDeflate {
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
	}

	if !cli.NoExtract {
		register.Archive("ar", ar.Detect, ar.Extract)
		register.Archive("cab", cab.Detect, cab.Extract)
		register.Archive("cpio", cpio.Detect, cpio.Extract)
		register.Archive("rar", rar.Detect, rar.Extract)
		register.Archive("rpm", rpm.Detect, rpm.Extract)
		register.Archive("szip", szip.Detect, szip.Extract)
		register.Archive("tar", tar.Detect, tar.Extract)
		register.Archive("xar", xar.Detect, xar.Extract)
		register.Archive("zip", zip.Detect, zip.Extract)
	}

	if !cli.NoConvert {
		register.Convert("evtx", evtx.Detect, evtx.Convert)
		register.Convert("journal", journal.Detect, journal.Convert)
		register.Convert("pe", pe.Detect, pe.Convert)
	}

	if !cli.Raw {
		register.Format("fox", fox.Detect, fox.Format)
		register.Format("json", json.Detect, json.Format)
	}

	cli.Heaps = heapset.New(args, &heapset.Options{
		Limit: &types.Limits{
			IsHead: cli.Head,
			IsTail: cli.Tail,
			Lines:  cli.Lines,
			Bytes:  cli.Bytes,
		},
		Filter: &types.Filters{
			Regex:  cli.Filter,
			Before: cli.Before,
			After:  cli.After,
		},
		Input:    cli.Input,
		Password: cli.Pass,
		Verbose:  cli.Verbose,
	})

	if cli.DryRun {
		for _, h := range cli.Heaps.Get() {
			_, _ = fmt.Fprintf(cli.Stdout, "%s\n", h.Name)
		}

		// exit early
		cli.Heaps.Discard()
		os.Exit(0)
	}

	return cli.Heaps
}

func (cli *Globals) Discard() {
	if len(cli.File) > 0 {
		_ = cli.Stdout.Close()
	}

	cli.Heaps.Discard()
}
