package cmd

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dlclark/regexp2/v2"
	"github.com/fatih/color"
	"go.foxforensics.eu/fox/v4/internal/pkg"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/fox/v4/internal/sys/heap"
	"go.foxforensics.eu/fox/v4/internal/sys/loader"
	"go.foxforensics.eu/fox/v4/internal/sys/memory"
	"go.foxforensics.eu/fox/v4/internal/sys/writer"
)

var Usage = strings.TrimSpace(`
Usage: fox [COMMAND] [FLAGS...] <PATHS...>

Commands:
   a, ad                   Show Active Directory infos
   c, cat                  Show file contents (default)
   s, str                  Show file contained strings
   i, info                 Show file infos and entropy
   t, time                 show file super timeline
   h, hash                 Show file hashes and checksums
   x, hunt                 Hunt critical system events

File flags:
  -I, --in=FILE            Read paths from file
  -O, --out=FILE           Write output to file (receipted)

Filter flags:
  -L, --limit=NUMBER       Filter using byte or line count
  -F, --find=REGEX         Filter using regular expression

Process flags:
  -T, --threads=CORES      Use cores for parallel threads
  -P, --password=TEXT      Use text to decrypt (7z, Rar, Zip)

Disable flags:
  -r, --raw[=LEVEL]        Don't process files (r/rr/rrr/rrrr)
  -q, --quiet              Don't print anything
      --no-pretty          Don't prettify the output
      --no-convert         Don't convert automatically
      --no-deflate         Don't deflate automatically
      --no-extract         Don't extract automatically
      --no-protect         Don't protect the loader
      --no-receipt         Don't write the receipt

Standard flags:
  -v, --verbose[=LEVEL]    Print more details (v/vv)
  -d, --dry-run            Print only the found files
      --version            Print the version number
      --help               Print this help message

Positional arguments:
  Globbing paths to open or '-' to read from STDIN

Example: Find occurrences in event logs
  $ fox -FWinlogon ./**/*.evtx

Example: Hunt down critical events
  $ fox hunt -u *.dd

Example: Show help on sub commands
  $ fox help hunt

Report bugs at: foxforensics.eu/issues
`)

type Globals struct {
	// file flags
	In  []byte `short:"I" long:"in" type:"filecontent"`
	Out string `short:"O" long:"out" xor:"out,quiet"`

	// filter flags
	Limit string `short:"L"`
	Find  string `short:"F"`

	// process flags
	Threads  int    `short:"T" default:"${cores}"`
	Password string `short:"P"`

	// disable flags
	Raw       int  `short:"r" type:"counter"`
	Quiet     bool `short:"q" xor:"out,quiet"`
	NoPretty  bool `long:"no-pretty"`
	NoConvert bool `long:"no-convert"`
	NoDeflate bool `long:"no-deflate"`
	NoExtract bool `long:"no-extract"`
	NoProtect bool `long:"no-protect"`
	NoReceipt bool `long:"no-receipt"`

	// standard flags
	Verbose int  `short:"v" type:"counter"`
	DryRun  bool `short:"d" long:"dry-run"`
	Help    bool

	// hidden
	Lexer string `hidden:""`
	Style string `hidden:""`

	// internal
	Context context.Context    `kong:"-"`
	Cancel  context.CancelFunc `kong:"-"`
	Writer  *writer.Writer     `kong:"-"`
	Regexp  *regexp2.Regexp    `kong:"-"`
	Loader  *loader.Loader     `kong:"-"`
	Query   pkg.Query          `kong:"-"`
	Paths   []string           `kong:"-"`
}

func (fox *Globals) Init(args []string, raw bool) (<-chan *heap.Heap, error) {
	var err error
	var lvl slog.Level

	switch fox.Verbose {
	case 0:
		lvl = slog.LevelWarn
	case 1:
		lvl = slog.LevelInfo
	default:
		lvl = slog.LevelDebug
	}

	slog.SetLogLoggerLevel(lvl)

	if raw {
		fox.NoConvert = true
	}

	if fox.Threads <= 0 {
		fox.Threads = 1 // must be at least one
	}

	// parse paths
	if len(fox.In) > 0 {
		fox.Paths = sys.ParseList(fox.In)
	}

	// redirect output
	if len(fox.Out) > 0 {
		f, err := sys.CreateFile(fox.Out)

		if err != nil {
			return nil, err
		}

		fox.NoPretty = true
		fox.Writer = writer.New(f)
	} else if fox.Quiet {
		fox.Writer = writer.New(io.Discard)
	} else {
		fox.Writer = writer.New(os.Stdout)
	}

	if len(fox.Find) > 0 {
		if fox.Regexp, err = regexp2.Compile(fox.Find); err != nil {
			return nil, errors.New("invalid regex syntax")
		}
	}

	err = pkg.SetQuery(&fox.Query, fox.Limit, fox.Regexp)

	if err != nil {
		return nil, err
	}

	if fox.Raw > 0 {
		fox.NoPretty = true
	}

	if fox.Raw > 1 {
		fox.NoConvert = true
	}

	if fox.Raw > 2 {
		fox.NoDeflate = true
		fox.NoExtract = true
	}

	if fox.Raw > 3 {
		fox.NoProtect = true
	}

	if len(fox.Lexer) > 0 {
		writer.Lexer = fox.Lexer
	}

	if len(fox.Style) > 0 {
		writer.Style = fox.Style
	}

	if !fox.NoDeflate {
		loader.RegisterDeflates()
	}

	if !fox.NoExtract {
		loader.RegisterExtracts()
	}

	if !fox.NoConvert {
		loader.RegisterConverts()
	}

	if !fox.NoPretty {
		loader.RegisterFormats()
	} else {
		color.NoColor = true // turn off color package
	}

	if fox.NoReceipt && len(fox.Out) > 0 {
		slog.Warn("receipts has been disabled!")
	}

	if fox.Password == "-" {
		if fox.Password, err = sys.Password(); err != nil {
			return nil, err
		}
	}

	fox.Loader = loader.New(&loader.Options{
		Query:    fox.Query,
		Protect:  !fox.NoProtect,
		Password: fox.Password,
	})

	// handle ctrl+c
	fox.Context, fox.Cancel = signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	heaps := fox.Loader.Load(fox.Context, args)

	if fox.DryRun {
		for h := range heaps {
			fox.Writer.WriteString(h.String())
		}

		// exit early
		fox.Exit(0)
	}

	return heaps, nil
}

func (fox *Globals) Exit(code int) {
	fox.Discard()
	os.Exit(code)
}

func (fox *Globals) Discard() {
	if fox.Cancel != nil {
		fox.Cancel()
	}

	if fox.Writer != nil {
		fox.Writer.Close(fox.Out, !fox.NoReceipt)
	}

	if fox.Loader != nil {
		fox.Loader.Exit()
	}

	memory.FreeAll()
}
