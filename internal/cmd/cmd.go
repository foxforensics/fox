package cmd

import (
	"log"
	"os"
	"regexp"

	"github.com/fatih/color"

	_zip "github.com/cuhsat/fox/v4/internal/pkg/data/archive/7z"

	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/ar"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/cab"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/cpio"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/iso"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/msi"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/rar"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/rpm"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/tar"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/xar"
	"github.com/cuhsat/fox/v4/internal/pkg/data/archive/zip"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/bin/elf"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/bin/ese"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/bin/lnk"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/bin/pe"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/bin/pf"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/log/evtx"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/log/fortinet"
	"github.com/cuhsat/fox/v4/internal/pkg/data/convert/log/journal"
	"github.com/cuhsat/fox/v4/internal/pkg/data/deflate/bgzf"
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
	"github.com/cuhsat/fox/v4/internal/pkg/data/format/jsonl"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/client"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heap"
	"github.com/cuhsat/fox/v4/internal/pkg/types/loader"
	"github.com/cuhsat/fox/v4/internal/pkg/types/register"
	"github.com/cuhsat/fox/v4/internal/pkg/types/smap"
)

type Globals struct {
	// file flags
	Paths string `short:"i" long:"in"`
	File  string `short:"o" long:"out" xor:"out,quiet"`

	// limit flags
	Head  bool   `short:"h" xor:"head,tail"`
	Tail  bool   `short:"t" xor:"head,tail"`
	Lines string `short:"l" xor:"lines,bytes"`
	Bytes string `short:"c" xor:"lines,bytes"`

	// filter flags
	Regex string `short:"e"`

	// special flags
	Password string `short:"p"`
	Parallel int    `short:"P" default:"${cores}"`

	// disable flags
	Raw       bool `short:"r"`
	Quiet     bool `short:"q" xor:"out,quiet"`
	NoPretty  bool `short:"N" long:"no-pretty"`
	NoStrict  bool `long:"no-strict"`
	NoDeflate bool `long:"no-deflate"`
	NoExtract bool `long:"no-extract"`
	NoConvert bool `long:"no-convert"`
	NoReceipt bool `long:"no-receipt"`

	// standard flags
	Verbose int  `short:"v" type:"counter"`
	DryRun  bool `short:"d" long:"dry-run"`
	Help    bool

	// hidden
	Lexer string `hidden:""`
	Style string `hidden:""`

	// internal
	Regexp *regexp.Regexp `kong:"-"`
	Loader *loader.Loader `kong:"-"`
	Filter *types.Filters `kong:"-"`
	Limit  *types.Limits  `kong:"-"`
}

func (cli *Globals) Load(args []string, raw bool) <-chan *heap.Heap {
	if raw {
		cli.NoConvert = true
	}

	if len(cli.File) > 0 {
		cli.NoPretty = true
	}

	if len(cli.Regex) > 0 {
		cli.Regexp = regexp.MustCompile(cli.Regex)
	}

	if len(cli.Lexer) > 0 {
		text.Lexer = cli.Lexer
	}

	if len(cli.Style) > 0 {
		text.Style = cli.Style
	}

	cli.Limit = types.NewLimits(
		cli.Head,
		cli.Tail,
		cli.Bytes,
		cli.Lines,
	)

	cli.Filter = &types.Filters{
		Regex: cli.Regexp,
	}

	if cli.Raw {
		cli.NoPretty = true
		cli.NoStrict = true
		cli.NoDeflate = true
		cli.NoExtract = true
		cli.NoConvert = true
		cli.NoReceipt = true
	}

	if cli.Parallel <= 0 {
		cli.Parallel = 1 // must be at least one
	}

	if !cli.NoDeflate {
		register.Deflate("bgzf", bgzf.Detect, bgzf.Deflate)
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
		register.Archive("7z", _zip.Detect, _zip.Extract)
		register.Archive("ar", ar.Detect, ar.Extract)
		register.Archive("cab", cab.Detect, cab.Extract)
		register.Archive("cpio", cpio.Detect, cpio.Extract)
		register.Archive("iso", iso.Detect, iso.Extract)
		register.Archive("msi", msi.Detect, msi.Extract)
		register.Archive("rar", rar.Detect, rar.Extract)
		register.Archive("rpm", rpm.Detect, rpm.Extract)
		register.Archive("tar", tar.Detect, tar.Extract)
		register.Archive("xar", xar.Detect, xar.Extract)
		register.Archive("zip", zip.Detect, zip.Extract)
	}

	if !cli.NoConvert {
		register.Convert("elf", elf.Detect, elf.Convert)
		register.Convert("ese", ese.Detect, ese.Convert)
		register.Convert("lnk", lnk.Detect, lnk.Convert)
		register.Convert("pe", pe.Detect, pe.Convert)
		register.Convert("pf", pf.Detect, pf.Convert)
		register.Convert("evtx", evtx.Detect, evtx.Convert)
		register.Convert("fortinet", fortinet.Detect, fortinet.Convert)
		register.Convert("journal", journal.Detect, journal.Convert)
	}

	if !cli.NoPretty {
		register.Format("json", json.Detect, json.Format)
		register.Format("jsonl", jsonl.Detect, jsonl.Format)
	} else {
		color.NoColor = true // turn off color package
	}

	if cli.NoReceipt && len(cli.File) > 0 {
		log.Println("warning: receipts has been disabled!")
	}

	cli.Loader = loader.New(&loader.Options{
		Limit:    cli.Limit,
		Filter:   cli.Filter,
		Paths:    cli.Paths,
		Password: cli.Password,
		Parallel: cli.Parallel,
		Verbose:  cli.Verbose,
		Strict:   !cli.NoStrict,
	})

	if cli.DryRun {
		for h := range cli.Loader.Load(args) {
			text.Write(h.Name)
		}

		// exit early
		cli.Exit(0)
	}

	client.Idle = cli.Parallel
	smap.Chunks = cli.Parallel

	return cli.Loader.Load(args)
}

func (cli *Globals) Exit(code int) {
	cli.Discard()
	os.Exit(code)
}

func (cli *Globals) Discard() {
	cli.Loader.Exit()
}
