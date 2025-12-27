package hunt

import (
	"fmt"
	"log"
	"strings"

	"github.com/alecthomas/kong"

	cli "github.com/cuhsat/fox/v4/internal/cmd"

	"github.com/cuhsat/fox/v4/internal/pkg/data/stream"
	"github.com/cuhsat/fox/v4/internal/pkg/data/stream/ecs"
	"github.com/cuhsat/fox/v4/internal/pkg/data/stream/hec"
	"github.com/cuhsat/fox/v4/internal/pkg/data/stream/raw"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
	"github.com/cuhsat/fox/v4/internal/pkg/types/hunter"
)

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
}

func (cmd *Hunt) Validate() error {
	if cmd.Hec && len(cmd.Auth) == 0 {
		log.Fatal("auth required")
	}

	return nil
}

func (cmd *Hunt) BeforeApply(_ *kong.Kong, _ kong.Vars) error {
	if len(cmd.Paths) == 0 {
		cmd.Paths = hunter.Paths
	}

	return nil
}

func (cmd *Hunt) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	if cmd.Logstash {
		cmd.Url = types.Logstash
		cmd.Ecs = true
	}

	if cmd.Splunk {
		cmd.Url = types.Splunk
		cmd.Hec = true
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

	var sa stream.Streamable
	var db *hunter.Database
	var fn text.Colored
	var tx int64
	var rx int64
	var s string

	cli.NoConvert = true // force

	ch := cli.Load(cmd.Paths)
	defer cli.Discard()

	if cli.Verbose > 0 {
		log.Println("hunt: started")
	}

	if cmd.Sqlite {
		db = hunter.NewDB(types.Database)

		if cli.Verbose > 0 {
			log.Printf("hunt: using database %s\n", db)
		}
	}

	if len(cmd.Url) > 0 {
		switch {
		case cmd.Hec:
			sa = hec.New(cmd.Url, cmd.Auth)
		case cmd.Ecs:
			sa = ecs.New(cmd.Url)
		default:
			sa = raw.New(cmd.Url)
		}

		if cli.Verbose > 0 {
			log.Printf("hunt: using schema %s\n", sa)
		}
	}

	cnt := 0

	for e := range hunter.New(&hunter.Options{
		Sort:       cmd.Sort,
		Extensions: cmd.Ext,
		Verbose:    cli.Verbose,
	}).Hunt(ch) {
		if cmd.All || e.Severity >= hunter.Level {
			// apply color
			switch {
			case cli.Filter != nil:
				fn = text.MarkMatchFunc(cli.Filter)
			case cmd.All && e.Severity >= hunter.Level:
				fn = text.Mark // mark event
			case cmd.All:
				fn = text.Hide // hide event
			default:
				fn = text.Term // reset terminal
			}

			// apply format
			switch {
			case cmd.Jsonl:
				s = fn(e.ToJSONL())
			case cmd.Json:
				s = fn(e.ToJSON())
			default:
				s = fn(e.ToCEF())
			}

			if cli.Filter != nil && !cli.Filter.MatchString(s) {
				continue // filter event
			}

			_, _ = fmt.Fprintln(cli.Stdout, s)

			// hook for database
			if db != nil {
				db.Persist(e)
			}

			// hook for stream
			if sa != nil {
				td, rd, _ := sa.Write(e)
				tx += td
				rx += rd
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

	if cli.Verbose > 2 && tx+rx > 0 {
		log.Println("hunt: stream tx:", text.Humanize(tx))
		log.Println("hunt: stream rx:", text.Humanize(rx))
	}

	return nil
}
