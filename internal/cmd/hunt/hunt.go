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

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/data/stream"
	"github.com/cuhsat/fox/v4/internal/pkg/data/stream/ecs"
	"github.com/cuhsat/fox/v4/internal/pkg/data/stream/hec"
	"github.com/cuhsat/fox/v4/internal/pkg/data/stream/raw"
	"github.com/cuhsat/fox/v4/internal/pkg/rules"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
	"github.com/cuhsat/fox/v4/internal/pkg/types/hunter"
	"github.com/cuhsat/fox/v4/internal/pkg/types/receipt"
)

var Usage = strings.TrimSpace(`
Hunts suspicious activities.

fox hunt [FLAGS...] [PATHS...]

Flags:
  -a, --all                shows logs with all severities
  -s, --sort               shows logs sorted by timestamp (slow)
  -j, --json               shows logs as JSON objects
  -J, --jsonl              shows logs as JSON lines
  -Q, --sqlite             saves logs to SQLite3 DB (very slow)

Rule Flags:
  -R, --rule=FILE          filters using a Sigma rule file (slow)

Stream Flags:
  -U, --url=SERVER         streams events to server address
  -A, --auth=TOKEN         streams events using auth token
  -E, --ecs                uses ECS schema for streaming
  -H, --hec                uses HEC schema for streaming

Aliases:
  -L, --logstash           alias for -E -Uhttp://localhost:8080
  -S, --splunk             alias for -H -Uhttp://localhost:8088/...

Examples:
  $ fox hunt -sv ./**/*.E01
`)

type Hunt struct {
	All    bool   `short:"a"`
	Sort   bool   `short:"s"`
	Json   bool   `short:"j" xor:"json,jsonl"`
	Jsonl  bool   `short:"J" xor:"json,jsonl"`
	Sqlite bool   `short:"Q"`
	Rule   string `short:"R" sep:"," type:"path"`

	// stream
	Url  string `short:"U"`
	Auth string `short:"A"`
	Ecs  bool   `short:"E" xor:"ecs,hec"`
	Hec  bool   `short:"H" xor:"ecs,hec"`

	// alias
	Logstash bool `short:"L" xor:"logstash,splunk"`
	Splunk   bool `short:"S" xor:"logstash,splunk"`

	// paths
	Paths []string `arg:"" type:"path" optional:""`

	// internal
	db   *event.Database `kong:"-"`
	net  stream.Streamer `kong:"-"`
	rule sigma.Rule      `kong:"-"`
}

func (cmd *Hunt) Validate() error {
	if cmd.Hec && len(cmd.Auth) == 0 {
		log.Fatalln("auth required")
	}

	return nil
}

func (cmd *Hunt) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	var err error

	rule := rules.Default

	if cmd.Sqlite {
		cmd.db = event.NewDB(hunter.Database)
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

	cli.NoConvert = true // forced

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
		log.Printf("hunt: using database %s\n", cmd.db)
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

		if !cmd.All && !res.Match {
			continue // not matched
		}

		line := cmd.format(e, cli.Regexp, res.Match)

		if cli.Regexp != nil && !cli.Regexp.MatchString(line) {
			continue // not matched afterward
		}

		_, _ = fmt.Fprintln(cli.Stdout, line)

		cmd.upsert(e)

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

func (cmd *Hunt) format(e *event.Event, re *regexp.Regexp, m bool) string {
	var fn text.Colored

	switch {
	case re != nil:
		fn = text.MarkMatchFunc(re)
	case cmd.All && m:
		fn = text.Mark // mark event
	case cmd.All:
		fn = text.Hide // hide event
	default:
		fn = text.Term // reset
	}

	switch {
	case cmd.Jsonl:
		return fn(e.ToJSONL())
	case cmd.Json:
		return fn(e.ToJSON())
	default:
		return fn(e.ToCEF())
	}
}

func (cmd *Hunt) upsert(e *event.Event) {
	if cmd.db != nil {
		err := cmd.db.Upsert(e)

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
	if cmd.db != nil && !cli.NoReceipt {
		err := receipt.Generate(hunter.Database)

		if err != nil {
			log.Println(err)
		}
	}
}
