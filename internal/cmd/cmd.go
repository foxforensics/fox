package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/fatih/color"

	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heapset"
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

func (cli *Globals) Bootstrap(args []string) *heapset.HeapSet {
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
		cli.NoConvert = true
		cli.NoDeflate = true
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
		Input:     cli.Input,
		Password:  cli.Pass,
		NoDeflate: cli.NoDeflate,
		NoConvert: cli.NoConvert,
		Verbose:   cli.Verbose,
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
