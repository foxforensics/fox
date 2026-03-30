package hunt

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/bradleyjkemp/sigma-go"
	"github.com/bradleyjkemp/sigma-go/evaluator"

	cli "go.foxforensics.dev/fox/v4/internal/cmd"

	"go.foxforensics.dev/fox/v4/internal/pkg/file/store"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/store/parquet"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/store/sqlite"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/stream"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/stream/ecs"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/stream/hec"
	"go.foxforensics.dev/fox/v4/internal/pkg/file/stream/raw"
	"go.foxforensics.dev/fox/v4/internal/pkg/rules"
	"go.foxforensics.dev/fox/v4/internal/pkg/text"
	"go.foxforensics.dev/fox/v4/internal/pkg/text/unique"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/event"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/hunter"
	"go.foxforensics.dev/fox/v4/internal/pkg/types/receipt"
)

var Usage = strings.TrimSpace(`
fox hunt [FLAGS...] [PATHS...]

Flags:
  -a, --all                Show logs with all severities
  -s, --sort               Show logs sorted by timestamp (slow)
  -u, --uniq               Show logs that are unique 
  -j, --json               Show logs as JSON objects
  -J, --jsonl              Show logs as JSON lines
  -P, --parquet            Save logs as Parquet (very fast)
  -S, --sqlite             Save logs as SQLite3 (very slow)

Hunter flags:
  -b, --block=SIZE         Block size for event carving

Filter flags:
  -R, --rule=FILE          Filter using Sigma Rules file (slow)
  -D, --dist=LENGTH        Filter using Levenshtein distance (slow)

Stream flags:
  -U, --url=SERVER         Stream events to server address
  -A, --auth=TOKEN         Stream events using auth token
  -E, --ecs                Use ECS schema for streaming
  -H, --hec                Use HEC schema for streaming

Aliases:
      --logstash           Alias for -EUhttp://localhost:8080
      --splunk             Alias for -HUhttp://localhost:8088/...

Example: Hunt down critical events
  $ fox hunt -u *.dd

Example: Save all events as Parquet
  $ fox hunt -aP *.evtx
`)

type Hunt struct {
	All     bool `short:"a"`
	Sort    bool `short:"s"`
	Uniq    bool `short:"u" xor:"uniq,dist"`
	Json    bool `short:"j" xor:"json,jsonl"`
	Jsonl   bool `short:"J" xor:"json,jsonl"`
	Parquet bool `short:"P" xor:"sqlite,parquet"`
	Sqlite  bool `short:"S" xor:"sqlite,parquet"`

	// hunter
	Block string `short:"b" default:"65536"`

	// filter
	Rule string  `short:"R"`
	Dist float64 `short:"D" xor:"uniq,dist"`

	// stream
	Url  string `short:"U"`
	Auth string `short:"A"`
	Ecs  bool   `short:"E" xor:"ecs,hec"`
	Hec  bool   `short:"H" xor:"ecs,hec"`

	// alias
	Logstash bool `xor:"logstash,splunk"`
	Splunk   bool `xor:"logstash,splunk"`

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

	if len(cmd.Block) > 0 {
		hunter.Block = int(text.Mechanize(cmd.Block))
	}

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

	if !cli.NoPretty {
		text.Title(cmd.Paths...)
	}

	ch := cli.Load(cmd.Paths, true)
	defer cli.Discard()

	defer cmd.discard(cli)

	if cli.Verbose > 0 {
		log.Println("hunt: started")
	}

	if cli.Verbose > 1 {
		log.Printf("hunt: using %d worker(s)\n", cli.Parallel)
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

	if !rules.IsSupported(&cmd.rule) {
		log.Println("warning: rule is not supported!")
	}

	var n int64

	var ctx = context.Background()
	var sig = evaluator.ForRule(cmd.rule)

	for e := range hunter.New(&hunter.Options{
		Sort:     cmd.Sort,
		Parallel: cli.Parallel,
		Verbose:  cli.Verbose,
	}).Hunt(ch) {
		m, err := sig.Matches(ctx, e.Fields)

		if err != nil {
			log.Println(err)
			continue // not successful
		}

		if cmd.uniq != nil && !cmd.uniq.IsUnique(e.String()) {
			continue // not unique
		}

		if !cmd.All && !m.Match {
			continue // not matched
		}

		if cmd.db == nil {
			text.Match(cmd.format(e), cli.Regexp)
		} else {
			err := cmd.db.Store(e)

			if err != nil {
				log.Println(err)
			}
		}

		if cmd.net != nil {
			err := cmd.net.Stream(e)

			if err != nil {
				log.Println(err)
			}
		}

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

func (cmd *Hunt) format(e *event.Event) string {
	switch {
	case cmd.Jsonl:
		return text.ColorizeAs(e.ToJSONL(), "json")
	case cmd.Json:
		return text.ColorizeAs(e.ToJSON(), "json")
	default:
		return text.MarkEvent(e.ToCEF())
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
