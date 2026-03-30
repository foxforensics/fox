package cmd

import (
	"log"
	"os"
	"regexp"

	"github.com/fatih/color"

	_zip "go.foxforensics.dev/fox/v4/internal/pkg/file/archive/7z"

	"go.foxforensics.dev/fox/v4/internal/pkg/file/archive/ar"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/archive/cab"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/archive/cpio"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/archive/iso"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/archive/msi"
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
	"go.foxforensics.dev/fox/v4/internal/pkg/file/convert/log/fortinet"
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
	"go.foxforensics.dev/fox/v4/internal/pkg/text"
	"go.foxforensics.dev/fox/v4/internal/pkg/types"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/client"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/heap"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/loader"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/register"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/smap"
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
	Parallel int    `short:"z" default:"${cores}"`

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
		register.Deflate("lznt1", lznt1.Detect, lznt1.Deflate)
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
		register.Format("xml", xml.Detect, xml.Format)
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
