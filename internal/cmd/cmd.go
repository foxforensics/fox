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
	"go.foxforensics.eu/checker/services/vt"

	"go.foxforensics.eu/fox/v4/internal/pkg/tables"
	"go.foxforensics.eu/fox/v4/internal/pkg/text"
	"go.foxforensics.eu/fox/v4/internal/pkg/types"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/client"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/heap"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/loader"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/smap"
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

func (cli *Globals) Init(args []string, raw bool) (<-chan *heap.Heap, error) {
	var err error

	if raw {
		cli.NoConvert = true
	}

	if len(cli.Out) > 0 {
		cli.NoPretty = true
	}

	if len(cli.Find) > 0 {
		if cli.Regexp, err = regexp2.Compile(cli.Find); err != nil {
			return nil, errors.New("invalid regex syntax")
		}
	}

	cli.Limits, err = types.NewLimits(cli.Limit)

	if err != nil {
		return nil, err
	}

	cli.Filters = &types.Filters{
		Regex: cli.Regexp,
	}

	if cli.Raw > 0 {
		cli.NoConvert = true
	}

	if cli.Raw > 1 {
		cli.NoDeflate = true
		cli.NoExtract = true
	}

	if cli.Raw > 2 {
		cli.NoStrict = true
	}

	if cli.Threads <= 0 {
		cli.Threads = 1 // must be at least one
	}

	if len(cli.Lexer) > 0 {
		text.Lexer = cli.Lexer
	}

	if len(cli.Style) > 0 {
		text.Style = cli.Style
	}

	if !cli.NoDeflate {
		loader.RegisterDeflates()
	}

	if !cli.NoExtract {
		loader.RegisterExtracts()
	}

	if !cli.NoConvert {
		loader.RegisterConverts()
	}

	if !cli.NoPretty {
		loader.RegisterFormats()
	} else {
		color.NoColor = true // turn off color package
	}

	if cli.NoReceipt && len(cli.Out) > 0 {
		slog.Warn("receipts has been disabled")
	}

	cli.Loader = loader.New(&loader.Options{
		Limits:   cli.Limits,
		Filters:  cli.Filters,
		Password: cli.Password,
		Threads:  cli.Threads,
		Verbose:  cli.Verbose,
		Strict:   !cli.NoStrict,
	})

	// handle CTRC+C
	cli.Context, cli.Cancel = signal.NotifyContext(
		context.Background(),
		os.Kill,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	vt.Key = cli.ApiKey
	smap.Threads = cli.Threads
	client.MaxIdle = cli.Threads
	tables.Threads = cli.Threads

	heaps := cli.Loader.Load(cli.Context, args)

	if cli.DryRun {
		for h := range heaps {
			text.Stdout.Write(h.Name)
		}

		// exit early
		cli.Exit(0)
	}

	return heaps, nil
}

func (cli *Globals) Exit(code int) {
	cli.Discard()
	os.Exit(code)
}

func (cli *Globals) Discard() {
	cli.Cancel()
	cli.Loader.Exit()
}
