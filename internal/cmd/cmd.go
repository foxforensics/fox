package cmd

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dlclark/regexp2/v2"
	"github.com/fatih/color"
	"go.foxforensics.eu/fox/v4/internal/net/client"
	"go.foxforensics.eu/fox/v4/internal/pkg/types"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/smap"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/tables"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/fox/v4/internal/sys/heap"
	"go.foxforensics.eu/fox/v4/internal/sys/loader"
	"go.foxforensics.eu/fox/v4/internal/sys/terminal"
)

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
	NoPretty  bool `short:"n" long:"no-pretty"`
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
	ApiKey string `hidden:"" long:"api-key"`
	Lexer  string `hidden:""`
	Style  string `hidden:""`

	// internal
	Context context.Context    `kong:"-"`
	Cancel  context.CancelFunc `kong:"-"`
	Regexp  *regexp2.Regexp    `kong:"-"`
	Loader  *loader.Loader     `kong:"-"`
	Filters *types.Filters     `kong:"-"`
	Limits  *types.Limits      `kong:"-"`
	Input   []string           `kong:"-"`
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

	if len(fox.Out) > 0 {
		fox.NoPretty = true
	}

	if len(fox.Find) > 0 {
		if fox.Regexp, err = regexp2.Compile(fox.Find); err != nil {
			return nil, errors.New("invalid regex syntax")
		}
	}

	fox.Limits, err = types.NewLimits(fox.Limit)

	if err != nil {
		return nil, err
	}

	fox.Filters = &types.Filters{
		Regex: fox.Regexp,
	}

	if fox.Raw > 0 {
		fox.NoConvert = true
	}

	if fox.Raw > 1 {
		fox.NoDeflate = true
		fox.NoExtract = true
	}

	if fox.Raw > 2 {
		fox.NoStrict = true
	}

	if fox.Threads <= 0 {
		fox.Threads = 1 // must be at least one
	}

	if len(fox.Lexer) > 0 {
		terminal.Lexer = fox.Lexer
	}

	if len(fox.Style) > 0 {
		terminal.Style = fox.Style
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

	fox.Loader = loader.New(&loader.Options{
		Limits:   fox.Limits,
		Filters:  fox.Filters,
		Password: fox.Password,
		Threads:  fox.Threads,
		Strict:   !fox.NoStrict,
	})

	// handle CTRC+C
	fox.Context, fox.Cancel = signal.NotifyContext(
		context.Background(),
		os.Kill,
		os.Interrupt,
		syscall.SIGTERM,
	)

	smap.Threads = fox.Threads
	client.MaxIdle = fox.Threads
	tables.Threads = fox.Threads

	heaps := fox.Loader.Load(fox.Context, args)

	if fox.DryRun {
		for h := range heaps {
			sys.Stdout.Write(h.Name)
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
	fox.Cancel()
	fox.Loader.Exit()
}
