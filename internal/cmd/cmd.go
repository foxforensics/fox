package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/fatih/color"

	"github.com/cuhsat/fox/v4/internal/pkg/data/stream"
	"github.com/cuhsat/fox/v4/internal/pkg/data/stream/ecs"
	"github.com/cuhsat/fox/v4/internal/pkg/data/stream/hec"
	"github.com/cuhsat/fox/v4/internal/pkg/data/stream/raw"
	"github.com/cuhsat/fox/v4/internal/pkg/hash"
	"github.com/cuhsat/fox/v4/internal/pkg/hunt"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/buffer"
	"github.com/cuhsat/fox/v4/internal/pkg/types/heapset"
)

type Hunt struct {
	All    bool     `short:"a"`
	Ext    int      `short:"x" type:"counter"`
	Sort   bool     `short:"s"`
	Json   bool     `short:"j" xor:"json,jsonl"`
	Jsonl  bool     `short:"J" xor:"json,jsonl"`
	Sqlite bool     `short:"D"`
	Paths  []string `arg:"" type:"path" optional:""`
}

type Hash struct {
	Algo  []string `short:"a" sep:"," default:"SHA256"`
	Find  []string `short:"F" sep:","`
	Paths []string `arg:"" type:"path" optional:""`
}

type Info struct {
	Min   float64  `short:"a" default:"0.0"`
	Max   float64  `short:"b" default:"1.0"`
	Paths []string `arg:"" name:"path" type:"path" optional:""`
}

type Text struct {
	Min   uint     `short:"a" default:"3"`
	Max   uint     `short:"b" default:"256"`
	Wtf   int      `short:"w" type:"counter"`
	First bool     `short:"1" and:"first,wtf"`
	Paths []string `arg:"" type:"path" optional:""`
}

type Hex struct {
	Mode  string   `short:"m" enum:"c,hd,xxd,raw" default:"raw"`
	Paths []string `arg:"" type:"path"`
}

type Cat struct {
	Paths []string `arg:"" type:"path" optional:""`
}

type Cli struct {
	// commands
	Hunt Hunt `cmd:"" aliases:"u"`
	Hash Hash `cmd:"" aliases:"h"`
	Info Info `cmd:"" aliases:"i,wc"`
	Text Text `cmd:"" aliases:"t,strings"`
	Hex  Hex  `cmd:"" aliases:"x"`
	Cat  Cat  `cmd:"" default:"withargs" aliases:"c,less"`

	// file limits
	Head  bool `short:"h" xor:"head,tail"`
	Tail  bool `short:"t" xor:"head,tail"`
	Lines uint `short:"n" xor:"lines,bytes"`
	Bytes uint `short:"c" xor:"lines,bytes"`

	// file loader
	Pass string `short:"p"`

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
	DryRun  bool `short:"d" long:"dry-run"`
	Verbose int  `short:"v" type:"counter"`

	// internal
	w  io.WriteCloser   `kong:"-"`
	re *regexp.Regexp   `kong:"-"`
	hs *heapset.HeapSet `kong:"-"`
}

func (cli *Cli) Bootstrap(args []string) *heapset.HeapSet {
	var sw io.Writer

	if len(cli.Regex) > 0 {
		cli.re = regexp.MustCompile(cli.Regex)
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

	if cli.Info.Min > cli.Info.Max {
		log.Fatal("invalid range")
	}

	if cli.Text.Min > cli.Text.Max {
		log.Fatal("invalid range")
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
		cli.w = stream.New(cli.File, sw)
	} else if cli.Quiet {
		log.SetOutput(io.Discard)
		cli.w, _ = os.Open(os.DevNull)
	} else {
		cli.w = os.Stdout
	}

	if len(cli.Hunt.Paths) == 0 {
		cli.Hunt.Paths = hunt.Paths
	}

	if cli.NoColor {
		color.NoColor = true
	}

	cli.hs = heapset.New(args, &heapset.Options{
		Limit: &types.Limits{
			IsHead: cli.Head,
			IsTail: cli.Tail,
			Lines:  cli.Lines,
			Bytes:  cli.Bytes,
		},
		Filter: &types.Filters{
			Regex:  cli.re,
			Before: cli.Before,
			After:  cli.After,
		},
		Password:  cli.Pass,
		NoDeflate: cli.NoDeflate,
		NoConvert: cli.NoConvert,
		Verbose:   cli.Verbose,
	})

	if cli.DryRun {
		for _, h := range cli.hs.Get() {
			_, _ = fmt.Fprintf(cli.w, "%s\n", h.Name)
		}

		// exit early
		cli.hs.ThrowAway()
		os.Exit(0)
	}

	return cli.hs
}

func (cli *Cli) ThrowAway() {
	if len(cli.File) > 0 {
		_ = cli.w.Close()
	}

	cli.hs.ThrowAway()
}

func (cmd *Hunt) Run(cli *Cli) error {
	var db *hunt.Database
	var fn text.Colored

	cli.NoConvert = true // force

	hs := cli.Bootstrap(cli.Hunt.Paths)
	defer cli.ThrowAway()

	if cli.Verbose > 0 {
		log.Println("hunt: started")
	}

	if cli.Hunt.Sqlite {
		db = hunt.NewDB(types.Database)

		if cli.Verbose > 0 {
			log.Printf("hunt: using %s\n", db)
		}
	}

	cnt := 0

	for e := range hunt.Hunt(hs, &hunt.Options{
		Sort:       cli.Hunt.Sort,
		Extensions: cli.Hunt.Ext,
		Verbose:    cli.Verbose,
	}) {
		if cli.Hunt.All || e.Severity >= hunt.Level {
			switch {
			case cli.Hunt.All && e.Severity >= hunt.Level:
				fn = text.Mark // mark event
			case cli.Hunt.All:
				fn = text.Hide // hide event
			default:
				fn = text.Term // reset terminal
			}

			switch {
			case cli.Hunt.Jsonl:
				_, _ = fmt.Fprintln(cli.w, fn(e.ToJSONL()))
			case cli.Hunt.Json:
				_, _ = fmt.Fprintln(cli.w, fn(e.ToJSON()))
			default:
				_, _ = fmt.Fprintln(cli.w, fn(e.ToCEF()))
			}

			if db != nil {
				db.Write(e)
			}

			cnt++
		}
	}

	if cli.Verbose > 0 {
		log.Println("hunt: finished")
	}

	if cli.Verbose > 1 {
		log.Printf("hunt: found %d events\n", cnt)
	}

	return nil
}

func (cmd *Hash) Run(cli *Cli) error {
	hs := cli.Bootstrap(cli.Hash.Paths)
	defer cli.ThrowAway()

	for _, algo := range cli.Hash.Algo {
		if len(cli.Hash.Algo) > 1 {
			_, _ = fmt.Fprintf(cli.w, "%s\n", text.Hide(text.Header(strings.ToUpper(algo))))
		}

		for _, h := range hs.Get() {
			sum, err := hash.Sum(algo, h.MMap())

			if err != nil {
				log.Println(err)
				continue
			}

			if len(cli.Hash.Find) == 0 || slices.Contains(cli.Hash.Find, sum) {
				_, _ = fmt.Fprintf(cli.w, "%s  %s\n", sum, text.Hide(h))
			}
		}
	}

	return nil
}

func (cmd *Info) Run(cli *Cli) error {
	hs := cli.Bootstrap(cli.Info.Paths)
	defer cli.ThrowAway()

	for _, h := range hs.Get() {
		if e, ok := h.Entropy(
			cli.Info.Min,
			cli.Info.Max,
		); ok {
			_, _ = fmt.Fprintf(cli.w, "%10dL %10dB  %.10fE  %s\n", h.Len(), len(h.MMap()), e, text.Hide(h.String()))
		}
	}

	return nil
}

func (cmd *Text) Run(cli *Cli) error {
	hs := cli.Bootstrap(cli.Text.Paths)
	defer cli.ThrowAway()

	for _, h := range hs.Get() {
		if hs.Len() > 1 && !cli.NoFile {
			_, _ = fmt.Fprintf(cli.w, "%s\n", text.Hide(text.Header(h.String())))
		}

		for s := range h.Strings(
			cli.Text.Min,
			cli.Text.Max,
			cli.Text.Wtf,
			cli.Text.First,
		) {
			if !cli.NoLine && cli.Text.Wtf > 0 {
				_, _ = fmt.Fprintf(cli.w, "%s  %s  %s\n", text.Hide(s.Off), s.Str, text.Hide(s.Cls))
			} else if !cli.NoLine {
				_, _ = fmt.Fprintf(cli.w, "%s  %s\n", text.Hide(s.Off), s.Str)
			} else {
				_, _ = fmt.Fprintf(cli.w, "%s\n", s.Str)
			}
		}
	}

	return nil
}

func (cmd *Hex) Run(cli *Cli) error {
	hs := cli.Bootstrap(cli.Hex.Paths)
	defer cli.ThrowAway()

	var tail uint

	if cli.Tail {
		tail = cli.Bytes
	}

	for _, h := range hs.Get() {
		if hs.Len() > 1 && !cli.NoFile {
			_, _ = fmt.Fprintf(cli.w, "%s\n", text.Hide(text.Header(h.String())))
		}

		lastHex, wasCut := "", false

		for l := range buffer.Hex(h, tail, cli.Hex.Mode).Lines {
			// cut similar lines for better readability
			if l.Hex == lastHex && cli.Hex.Mode != types.Raw {
				if !wasCut {
					wasCut = true
					_, _ = fmt.Fprintln(cli.w, text.Hide("*"))
				}
				continue
			}

			switch cli.Hex.Mode {
			case types.Canonical:
				_, _ = fmt.Fprintf(cli.w, "%s  %s%s\n", text.Hide(l.Nr), l.Hex, text.Hide(l.Str))
			case types.Hexdump:
				_, _ = fmt.Fprintf(cli.w, "%s %s\n", text.Hide(l.Nr), l.Hex)
			case types.Xxd:
				_, _ = fmt.Fprintf(cli.w, "%s %s %-16s\n", text.Hide(l.Nr), l.Hex, text.Hide(l.Str))
			case types.Raw:
				_, _ = fmt.Fprintf(cli.w, "%s\n", l.Hex)
			}

			lastHex, wasCut = l.Hex, false
		}
	}

	return nil
}

func (cmd *Cat) Run(cli *Cli) error {
	hs := cli.Bootstrap(cli.Cat.Paths)
	defer cli.ThrowAway()

	for _, h := range hs.Get() {
		if hs.Len() > 1 && !cli.NoFile {
			_, _ = fmt.Fprintf(cli.w, "%s\n", text.Hide(text.Header(h.String())))
		}

		for l := range buffer.Text(h, 2).Lines {
			s := l.Str

			if cli.re != nil && l.Nr != buffer.Sep {
				s = text.MarkMatch(s, cli.re)
			}

			if !cli.NoLine && l.Nr == buffer.Sep {
				_, _ = fmt.Fprintf(cli.w, "%s\n", text.Hide(buffer.Sep))
			} else if !cli.NoLine {
				_, _ = fmt.Fprintf(cli.w, "%s %s\n", text.Hide(l.Nr), s)
			} else {
				_, _ = fmt.Fprintf(cli.w, "%s\n", s)
			}
		}
	}

	return nil
}
