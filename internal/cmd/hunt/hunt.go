package hunt

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/alecthomas/kong"
	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/data/stream"
	"github.com/cuhsat/fox/v4/internal/pkg/data/stream/ecs"
	"github.com/cuhsat/fox/v4/internal/pkg/data/stream/hec"
	"github.com/cuhsat/fox/v4/internal/pkg/data/stream/raw"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
	"github.com/cuhsat/fox/v4/internal/pkg/types/hunter"
)

const Limit = 8

var Usage = strings.TrimSpace(`
Hunt suspicious activities.

fox hunt [FLAGS ...] [PATHS ...]

Flags:
  -a, --all                show logs with all severities
  -x, --ext                show logs with all extensions (slow)
  -s, --sort               show logs sorted by timestamp (slow)
  -j, --json               show logs as JSON objects
  -J, --jsonl              show logs as JSON lines
  -D, --sqlite             save logs to SQLite3 DB (very slow)

Worker:
  -P, --pool=SIZE          use worker pool size (default: CPUs)

Stream:
  -u, --url=SERVER         stream events to server address
  -T, --auth=TOKEN         stream events using auth token
  -E, --ecs                use ECS schema for streaming
  -H, --hec                use HEC schema for streaming

Alias:
  -L, --logstash           alias for -E -uhttp://localhost:8080
  -S, --splunk             alias for -H -uhttp://localhost:8088/...

Example:
  $ fox hunt -sxv ./**/*.dd
`)

type Hunt struct {
	All    bool `short:"a"`
	Ext    int  `short:"x" type:"counter"`
	Sort   bool `short:"s"`
	Json   bool `short:"j" xor:"json,jsonl"`
	Jsonl  bool `short:"J" xor:"json,jsonl"`
	Sqlite bool `short:"D"`

	// worker
	Pool int `short:"P" default:"${cpus}"`

	// stream
	Url  string `short:"u"`
	Auth string `short:"T"`
	Ecs  bool   `short:"E" xor:"ecs,hec"`
	Hec  bool   `short:"H" xor:"ecs,hec" and:"hec"`

	// alias
	Logstash bool `short:"L" xor:"logstash,splunk"`
	Splunk   bool `short:"S" xor:"logstash,splunk"`

	// paths
	Paths []string `arg:"" type:"path" optional:""`

	// internal
	db  *event.Database `kong:"-"`
	net stream.Streamer `kong:"-"`
}

func (cmd *Hunt) Validate() error {
	if cmd.Pool < 1 {
		log.Fatalln("pool to small")
	}

	if cmd.Hec && len(cmd.Auth) == 0 {
		log.Fatalln("auth required")
	}

	return nil
}

func (cmd *Hunt) BeforeApply(_ *kong.Kong, _ kong.Vars) error {
	if len(cmd.Paths) == 0 {
		cmd.Paths = hunter.Local
	}

	return nil
}

func (cmd *Hunt) AfterApply(_ *kong.Kong, _ kong.Vars) error {
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

	// extensions must be activated
	if cmd.Ecs || cmd.Hec {
		cmd.Ext = 3
	}

	return nil
}

func (cmd *Hunt) Run(cli *cli.Globals) error {
	if cli.Help {
		fmt.Print(Usage)
		return nil
	}

	cli.NoConvert = true // forced

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	if cli.Verbose > 0 {
		log.Println("hunt: started")
	}

	if cli.Verbose > 1 {
		log.Printf("hunt: using %d worker(s)\n", cmd.Pool)
	}

	if cli.Verbose > 1 && cmd.db != nil {
		log.Printf("hunt: using database %s\n", cmd.db)
	}

	if cli.Verbose > 1 && cmd.net != nil {
		log.Printf("hunt: streaming as %s\n", cmd.net)
	}

	var n int64

	for e := range hunter.New(&hunter.Options{
		Sort:       cmd.Sort,
		Extensions: cmd.Ext,
		Pool:       cmd.Pool,
		Verbose:    cli.Verbose,
	}).Hunt(ch) {
		if !cmd.All && e.Severity < Limit {
			continue // not severe enough
		}

		line := cmd.format(e, cli.Filter)

		if cli.Filter != nil && !cli.Filter.MatchString(line) {
			continue // not matched
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

func (cmd *Hunt) format(e *event.Event, re *regexp.Regexp) string {
	var fn text.Colored

	switch {
	case re != nil:
		fn = text.MarkMatchFunc(re)
	case cmd.All && e.Severity >= Limit:
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
