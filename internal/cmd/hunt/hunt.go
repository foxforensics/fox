package hunt

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/bradleyjkemp/sigma-go"
	"github.com/bradleyjkemp/sigma-go/evaluator"

	cli "foxhunt.dev/fox/internal/cmd"

	"foxhunt.dev/fox/internal/pkg/data/reader/ewf"
	"foxhunt.dev/fox/internal/pkg/data/reader/vhdx"
	"foxhunt.dev/fox/internal/pkg/data/reader/vmdk"
	"foxhunt.dev/fox/internal/pkg/data/store"
	"foxhunt.dev/fox/internal/pkg/data/store/parquet"
	"foxhunt.dev/fox/internal/pkg/data/store/sqlite"
	"foxhunt.dev/fox/internal/pkg/data/stream"
	"foxhunt.dev/fox/internal/pkg/data/stream/ecs"
	"foxhunt.dev/fox/internal/pkg/data/stream/hec"
	"foxhunt.dev/fox/internal/pkg/data/stream/raw"
	"foxhunt.dev/fox/internal/pkg/rules"
	"foxhunt.dev/fox/internal/pkg/text"
	"foxhunt.dev/fox/internal/pkg/text/unique"
	"foxhunt.dev/fox/internal/pkg/types/event"
	"foxhunt.dev/fox/internal/pkg/types/hunter"
	"foxhunt.dev/fox/internal/pkg/types/receipt"
	"foxhunt.dev/fox/internal/pkg/types/register"
)

var Usage = strings.TrimSpace(`
Hunts suspicious events.

fox hunt [FLAGS...] [PATHS...]

Flags:
  -a, --all                shows logs with all severities
  -s, --sort               shows logs sorted by timestamp (slow)
  -u, --uniq               shows logs that are unique 
  -j, --json               shows logs as JSON objects
  -J, --jsonl              shows logs as JSON lines
  -P, --parquet            saves logs as Parquet (very fast)
  -Q, --sqlite             saves logs as SQLite3 (very slow)

Hunter flags:
  -b, --block=SIZE         block size for event carving

Filter flags:
  -R, --rule=FILE          filters using Sigma Rules file (slow)
  -D, --dist=LENGTH        filters using Levenshtein distance (slow)

Stream flags:
  -U, --url=SERVER         streams events to server address
  -A, --auth=TOKEN         streams events using auth token
  -E, --ecs                uses ECS schema for streaming
  -H, --hec                uses HEC schema for streaming

Aliases:
  -L, --logstash           alias for -E -Uhttp://localhost:8080
  -S, --splunk             alias for -H -Uhttp://localhost:8088/...

Examples:
  $ fox hunt -sv ./**/*.dd
`)

type Hunt struct {
	All     bool `short:"a"`
	Sort    bool `short:"s"`
	Uniq    bool `short:"u" xor:"uniq,dist"`
	Json    bool `short:"j" xor:"json,jsonl"`
	Jsonl   bool `short:"J" xor:"json,jsonl"`
	Sqlite  bool `short:"Q" xor:"sqlite,parquet"`
	Parquet bool `short:"P" xor:"sqlite,parquet"`

	// hunter
	Block uint `short:"b" default:"65536"`

	// filter
	Rule string  `short:"R"`
	Dist float64 `short:"D" xor:"uniq,dist"`

	// stream
	Url  string `short:"U"`
	Auth string `short:"A"`
	Ecs  bool   `short:"E" xor:"ecs,hec"`
	Hec  bool   `short:"H" xor:"ecs,hec"`

	// alias
	Logstash bool `short:"L" xor:"logstash,splunk"`
	Splunk   bool `short:"S" xor:"logstash,splunk"`

	// paths
	Paths []string `arg:"" optional:""`

	// internal
	db   store.Store     `kong:"-"`
	net  stream.Streamer `kong:"-"`
	rule sigma.Rule      `kong:"-"`
	uniq unique.Unique   `kong:"-"`
}

func (cmd *Hunt) Validate() error {
	if cmd.Hec && len(cmd.Auth) == 0 {
		log.Fatalln("auth required")
	}

	return nil
}

func (cmd *Hunt) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	var err error

	hunter.Block = int(cmd.Block)

	switch {
	case cmd.Uniq:
		cmd.uniq = unique.ByHash()
	case cmd.Dist > 0:
		cmd.uniq = unique.ByDistance(cmd.Dist)
	}

	if cmd.Sqlite {
		cmd.db = sqlite.New(hunter.Storage)
	}

	if cmd.Parquet {
		cmd.db = parquet.New(hunter.Storage)
	}

	if cmd.Logstash {
		cmd.Url = ecs.LocalHost
		cmd.Ecs = true
	}

	if cmd.Splunk {
		cmd.Url = hec.LocalHost
		cmd.Hec = true
	}

	if len(cmd.Url) > 0 {
		switch {
		case cmd.Hec:
			cmd.net = hec.New(cmd.Url, cmd.Auth)
		case cmd.Ecs:
			cmd.net = ecs.New(cmd.Url)
		default:
			cmd.net = raw.New(cmd.Url)
		}
	}

	rule := rules.Default

	if len(cmd.Rule) > 0 {
		rule, err = os.ReadFile(cmd.Rule)

		if err != nil {
			log.Fatalln(err)
		}
	}

	cmd.rule, err = sigma.ParseRule(rule)

	if err != nil {
		log.Fatalln(err)
	}

	return nil
}

func (cmd *Hunt) Run(cli *cli.Globals) error {
	if len(cmd.Paths)+len(cli.Paths) == 0 {
		cmd.Paths = hunter.Local
	}

	if cmd.Dist > 0 {
		cli.NoSyntax = true
	}

	if !cli.Raw {
		register.Reader("ewf", ewf.Detect, ewf.Reader)
		register.Reader("vhdx", vhdx.Detect, vhdx.Reader)
		register.Reader("vmdk", vmdk.Detect, vmdk.Reader)
	}

	if !cli.NoStrict {
		cli.NoExtract = true // forced
		cli.NoDeflate = true // forced
		cli.NoConvert = true // forced
	}

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	defer cmd.discard(cli)

	if cli.Verbose > 0 {
		log.Println("hunt: started")
	}

	if cli.Verbose > 1 {
		log.Printf("hunt: using %d worker(s)\n", cli.Threads)
	}

	if cli.Verbose > 1 && cmd.db != nil {
		log.Printf("hunt: using store %s\n", cmd.db)
	}

	if cli.Verbose > 1 {
		log.Printf("hunt: using rule \"%s\"\n", cmd.rule.Title)
	}

	if cli.Verbose > 1 && cmd.net != nil {
		log.Printf("hunt: streaming as %s\n", cmd.net)
	}

	if !rules.IsSupported(&cmd.rule) && !cli.NoWarnings {
		log.Println("warning: rule is not supported!")
	}

	var n int64

	var ctx = context.Background()
	var sig = evaluator.ForRule(cmd.rule)

	for e := range hunter.New(&hunter.Options{
		Sort:     cmd.Sort,
		Parallel: cli.Threads,
		Verbose:  cli.Verbose,
	}).Hunt(ch) {
		res, err := sig.Matches(ctx, e.Fields)

		if err != nil {
			log.Println(err)
			continue // not successful
		}

		if cmd.uniq != nil && !cmd.uniq.IsUnique(e.String()) {
			continue // not unique
		}

		if !cmd.All && !res.Match {
			continue // not matched
		}

		line := cmd.format(e, cli.Regexp)

		if cli.Regexp != nil && !cli.Regexp.MatchString(line) {
			continue // not matched afterward
		}

		if cmd.db == nil {
			_, _ = fmt.Fprintln(cli.Stdout, line)
		}

		cmd.store(e)

		cmd.stream(e)

		n++
	}

	if cli.Verbose > 0 {
		log.Println("hunt: finished")
	}

	if cli.Verbose > 1 {
		log.Printf("hunt: found %d event(s)\n", n)
	}

	return nil
}

func (cmd *Hunt) format(e *event.Event, re *regexp.Regexp) string {
	var line string

	switch {
	case cmd.Jsonl:
		line = text.ColorizeStringAs(e.ToJSONL(), "json")
	case cmd.Json:
		line = text.ColorizeStringAs(e.ToJSON(), "json")
	default:
		line = e.ToCEF()
	}

	if re != nil {
		line = text.MarkMatch(line, re)
	}

	return line
}

func (cmd *Hunt) store(e *event.Event) {
	if cmd.db != nil {
		err := cmd.db.Store(e)

		if err != nil {
			log.Println(err)
		}
	}
}

func (cmd *Hunt) stream(e *event.Event) {
	if cmd.net != nil {
		err := cmd.net.Stream(e)

		if err != nil {
			log.Println(err)
		}
	}
}

func (cmd *Hunt) discard(cli *cli.Globals) {
	if cmd.db == nil {
		return
	}

	err := cmd.db.Close()

	if err != nil {
		log.Println(err)
	}

	if cli.NoReceipt {
		return
	}

	err = receipt.Generate(cmd.db.String())

	if err != nil {
		log.Println(err)
	}
}
