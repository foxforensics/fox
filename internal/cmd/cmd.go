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
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/bin/elf"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/bin/lnk"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/bin/pe"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/bin/pf"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/log/evtx"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/log/journal"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/br"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/bzip2"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/gzip"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/kanzi"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/lz4"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/lzfse"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/lzip"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/lzo"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/lzw"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/minlz"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/s2"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/snappy"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/xz"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/zlib"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/zstd"
	"github.com/cuhsat/fox/v4/internal/pkg/data/format/json"
	"github.com/cuhsat/fox/v4/internal/pkg/data/image/ewf"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heap"
	"github.com/cuhsat/fox/v4/internal/pkg/types/loader"
	"github.com/cuhsat/fox/v4/internal/pkg/types/receipt"
	"github.com/cuhsat/fox/v4/internal/pkg/types/register"
	"github.com/cuhsat/fox/v4/internal/pkg/types/smap"
)

type Globals struct {
	// file limits
	Head  bool `short:"h" xor:"head,tail"`
	Tail  bool `short:"t" xor:"head,tail"`
	Lines uint `short:"n" xor:"lines,bytes"`
	Bytes uint `short:"c" xor:"lines,bytes"`

	// file loader
	Pass  string `short:"p"`
	File  string `short:"f" type:"path"`
	Input string `short:"i"`

	// line output
	Output string `short:"o" xor:"output,quiet"`

	// line filter
	Regex   string `short:"e"`
	Context uint   `short:"C"`
	Before  uint   `short:"B"`
	After   uint   `short:"A"`

	// parallel
	Parallel int `short:"P" default:"${cores}"`

	// disable
	Raw        bool `short:"r"`
	Quiet      bool `short:"q" xor:"output,quiet"`
	NoFile     bool `long:"no-file"`
	NoLine     bool `long:"no-line"`
	NoColor    bool `long:"no-color"`
	NoPretty   bool `long:"no-pretty"`
	NoMapping  bool `long:"no-mapping"`
	NoDeflate  bool `long:"no-deflate"`
	NoExtract  bool `long:"no-extract"`
	NoConvert  bool `long:"no-convert"`
	NoReceipt  bool `long:"no-receipt"`
	NoWarnings bool `long:"no-warnings"`

	// standard
	Help    bool
	Less    bool `short:"l"`
	DryRun  bool `short:"d" long:"dry-run"`
	Verbose int  `short:"v" type:"counter"`

	// bootstrap
	Stdout io.WriteCloser `kong:"-"`
	Regexp *regexp.Regexp `kong:"-"`
	Loader *loader.Loader `kong:"-"`
	Filter *types.Filters `kong:"-"`
	Limit  *types.Limits  `kong:"-"`
}

func (cli *Globals) Load(args []string) <-chan *heap.Heap {
	var err error

	switch {
	// file with receipt
	case len(cli.Output) > 0:
		cli.NoColor = true
		cli.Stdout, err = os.OpenFile(cli.Output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)

		if err != nil {
			log.Fatalln(err)
		}

	// disabled
	case cli.Quiet:
		log.SetOutput(io.Discard)
		cli.Stdout, _ = os.Open(os.DevNull)

	// stdout
	default:
		cli.Stdout = os.Stdout
	}

	if len(cli.Regex) > 0 {
		cli.Regexp = regexp.MustCompile(cli.Regex)
	}

	if cli.Context > 0 {
		cli.Before = cli.Context
		cli.After = cli.Context
	}

	cli.Limit = &types.Limits{
		IsHead: cli.Head,
		IsTail: cli.Tail,
		Lines:  cli.Lines,
		Bytes:  cli.Bytes,
	}

	cli.Filter = &types.Filters{
		Regex:  cli.Regexp,
		Before: cli.Before,
		After:  cli.After,
	}

	if cli.Raw {
		cli.NoFile = true
		cli.NoLine = true
		cli.NoColor = true
		cli.NoPretty = true
		cli.NoMapping = true
		cli.NoDeflate = true
		cli.NoExtract = true
		cli.NoConvert = true
		cli.NoReceipt = true
		cli.NoWarnings = true
	}

	if cli.Parallel <= 0 {
		cli.Parallel = 1 // must be at least one
	}

	if cli.NoColor {
		color.NoColor = true // turn off color package
	}

	if cli.NoReceipt && !cli.NoWarnings {
		log.Println("warning: receipts has been disabled!")
	}

	if !cli.NoDeflate {
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
		register.Convert("elf", elf.Detect, elf.Convert)
		register.Convert("lnk", lnk.Detect, lnk.Convert)
		register.Convert("pe", pe.Detect, pe.Convert)
		register.Convert("pf", pf.Detect, pf.Convert)
		register.Convert("evtx", evtx.Detect, evtx.Convert)
		register.Convert("journal", journal.Detect, journal.Convert)
	}

	if !cli.NoPretty {
		register.Format("json", json.Detect, json.Format)
	}

	if !cli.Raw {
		register.Image("ewf", ewf.Detect, ewf.Ingest)
	}

	cli.Loader = loader.New(&loader.Options{
		Limit:    cli.Limit,
		Filter:   cli.Filter,
		File:     cli.File,
		Input:    cli.Input,
		Password: cli.Pass,
		Parallel: cli.Parallel,
		Verbose:  cli.Verbose,
		Warning:  !cli.NoWarnings,
	})

	if cli.DryRun {
		for h := range cli.Loader.Load(args) {
			_, _ = fmt.Fprintf(cli.Stdout, "%s\n", h.Name)
		}

		// exit early
		cli.Loader.Exit()
		os.Exit(0)
	}

	smap.Chunks = cli.Parallel

	return cli.Loader.Load(args)
}

func (cli *Globals) Discard() {
	if len(cli.Output) > 0 {
		_ = cli.Stdout.Close()

		if !cli.NoReceipt {
			err := receipt.Generate(cli.Output)

			if err != nil {
				log.Println(err)
			}
		}
	}

	cli.Loader.Exit()
}
