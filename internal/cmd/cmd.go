package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/fatih/color"

	"github.com/cuhsat/fox/v4/internal/pkg/data/stream"
	"github.com/cuhsat/fox/v4/internal/pkg/data/stream/ecs"
	"github.com/cuhsat/fox/v4/internal/pkg/data/stream/hec"
	"github.com/cuhsat/fox/v4/internal/pkg/data/stream/raw"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heapset"
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

	// line filter
	Regex   string `short:"e"`
	Context uint   `short:"C"`
	Before  uint   `short:"B"`
	After   uint   `short:"A"`

	// data stream
	File string `short:"f"`
	Url  string `short:"u"`
	Auth string `short:"T"`
	Ecs  bool   `short:"E" xor:"ecs,hec"`
	Hec  bool   `short:"H" xor:"ecs,hec" and:"hec,auth"`

	// disable
	Raw       bool `short:"r"`
	Quiet     bool `short:"q"`
	NoFile    bool `long:"no-file"`
	NoLine    bool `long:"no-line"`
	NoColor   bool `long:"no-color"`
	NoDeflate bool `long:"no-deflate"`
	NoConvert bool `long:"no-convert"`

	// aliases
	Logstash bool `short:"L" xor:"logstash,splunk"`
	Splunk   bool `short:"S" xor:"logstash,splunk"`

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
	var sw io.Writer

	if len(cli.Regex) > 0 {
		cli.Filter = regexp.MustCompile(cli.Regex)
	}

	if len(cli.Url) > 0 {
		switch {
		case cli.Hec:
			sw = hec.New(cli.Url, cli.Auth)
		case cli.Ecs:
			sw = ecs.New(cli.Url)
		default:
			sw = raw.New(cli.Url)
		}
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

	if cli.Logstash {
		cli.Url = types.Logstash
		cli.Ecs = true
	}

	if cli.Splunk {
		cli.Url = types.Splunk
		cli.Hec = true
	}

	if len(cli.File)+len(cli.Url) > 0 {
		cli.NoColor = true
		cli.Stdout = stream.New(cli.File, sw)
	} else if cli.Quiet {
		log.SetOutput(io.Discard)
		cli.Stdout, _ = os.Open(os.DevNull)
	} else {
		cli.Stdout = os.Stdout
	}

	if cli.NoColor {
		color.NoColor = true
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
		cli.Heaps.ThrowAway()
		os.Exit(0)
	}

	return cli.Heaps
}

func (cli *Globals) ThrowAway() {
	if len(cli.File) > 0 {
		_ = cli.Stdout.Close()
	}

	cli.Heaps.ThrowAway()
}
