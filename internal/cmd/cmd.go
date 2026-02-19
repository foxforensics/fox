package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	_color "github.com/fatih/color"

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
	"foxhunt.dev/fox/internal/pkg/text"
	"foxhunt.dev/fox/internal/pkg/types"
	"foxhunt.dev/fox/internal/pkg/types/client"
	"foxhunt.dev/fox/internal/pkg/types/heap"
	"foxhunt.dev/fox/internal/pkg/types/loader"
	"foxhunt.dev/fox/internal/pkg/types/receipt"
	"foxhunt.dev/fox/internal/pkg/types/register"
	"foxhunt.dev/fox/internal/pkg/types/smap"
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

	// archive flags
	Password string `short:"p"`

	// profile flags
	Threads int `short:"T" default:"${cores}"`

	// disable flags
	Raw        bool `short:"r"`
	Quiet      bool `short:"q" xor:"out,quiet"`
	NoFile     bool `long:"no-file"`
	NoLine     bool `long:"no-line"`
	NoColor    bool `long:"no-color"`
	NoSyntax   bool `long:"no-syntax"`
	NoPretty   bool `long:"no-pretty"`
	NoStrict   bool `long:"no-strict"`
	NoMapping  bool `long:"no-mapping"`
	NoDeflate  bool `long:"no-deflate"`
	NoExtract  bool `long:"no-extract"`
	NoConvert  bool `long:"no-convert"`
	NoReceipt  bool `long:"no-receipt"`
	NoWarnings bool `long:"no-warnings"`

	// standard flags
	Help    bool
	DryRun  bool `short:"d" long:"dry-run"`
	Verbose int  `short:"v" type:"counter"`

	// internal
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
	case len(cli.File) > 0:
		cli.NoColor = true
		cli.Stdout, err = os.OpenFile(cli.File, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)

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
		cli.NoSyntax = true
		cli.Regexp = regexp.MustCompile(cli.Regex)
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
		cli.NoFile = true
		cli.NoLine = true
		cli.NoColor = true
		cli.NoSyntax = true
		cli.NoPretty = true
		cli.NoStrict = true
		cli.NoMapping = true
		cli.NoDeflate = true
		cli.NoExtract = true
		cli.NoConvert = true
		cli.NoReceipt = true
		cli.NoWarnings = true
	}

	if cli.Threads <= 0 {
		cli.Threads = 1 // must be at least one
	}

	if cli.NoColor {
		cli.NoSyntax = true
		_color.NoColor = true // turn off color package
	}

	if cli.NoSyntax {
		text.NoSyntax = true // turn off syntax highlighting
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
		register.Archive("iso", iso.Detect, iso.Extract)
		register.Archive("rar", rar.Detect, rar.Extract)
		register.Archive("rpm", rpm.Detect, rpm.Extract)
		register.Archive("7z", _zip.Detect, _zip.Extract)
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
	}

	cli.Loader = loader.New(&loader.Options{
		Limit:    cli.Limit,
		Filter:   cli.Filter,
		Paths:    cli.Paths,
		Password: cli.Password,
		Parallel: cli.Threads,
		Verbose:  cli.Verbose,
		Strict:   !cli.NoStrict,
		Warnings: !cli.NoWarnings,
	})

	if cli.DryRun {
		for h := range cli.Loader.Load(args) {
			_, _ = fmt.Fprintf(cli.Stdout, "%s\n", h.Name)
		}

		// exit early
		cli.Loader.Exit()
		os.Exit(0)
	}

	client.Idle = cli.Threads
	smap.Chunks = cli.Threads

	return cli.Loader.Load(args)
}

func (cli *Globals) Exit(code int) {
	cli.Discard()

	os.Exit(code)
}

func (cli *Globals) Discard() {
	if len(cli.File) > 0 {
		_ = cli.Stdout.Close()

		if !cli.NoReceipt {
			err := receipt.Generate(cli.File)

			if err != nil {
				log.Println(err)
			}
		}
	}

	cli.Loader.Exit()
}
