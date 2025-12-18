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
	"github.com/cuhsat/fox/v4/internal/pkg/hunt"
	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
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
  -D, --sqlite             save logs to SQLite3 DB

Streams:
  -u, --url=SERVER         stream events to server address
  -T, --auth=TOKEN         stream events using auth token
  -E, --ecs                use ECS schema for streaming
  -H, --hec                use HEC schema for streaming

Aliases:
  -L, --logstash           alias for -E -uhttp://localhost:8080
  -S, --splunk             alias for -H -uhttp://localhost:8088/...
`)

type Hunt struct {
	All    bool `short:"a"`
	Ext    int  `short:"x" type:"counter"`
	Sort   bool `short:"s"`
	Json   bool `short:"j" xor:"json,jsonl"`
	Jsonl  bool `short:"J" xor:"json,jsonl"`
	Sqlite bool `short:"D"`

	// streams
	Url  string `short:"u"`
	Auth string `short:"T"`
	Ecs  bool   `short:"E" xor:"ecs,hec"`
	Hec  bool   `short:"H" xor:"ecs,hec" and:"hec,auth"`

	// aliases
	Logstash bool `short:"L" xor:"logstash,splunk"`
	Splunk   bool `short:"S" xor:"logstash,splunk"`

	// paths
	Paths []string `arg:"" type:"path" optional:""`
}

func (cmd *Hunt) BeforeApply(_ *kong.Kong, _ kong.Vars) error {
	if len(cmd.Paths) == 0 {
		cmd.Paths = hunt.Paths
	}

	if cmd.Logstash {
		cmd.Url = types.Logstash
		cmd.Ecs = true
	}

	if cmd.Splunk {
		cmd.Url = types.Splunk
		cmd.Hec = true
	}

	return nil
}

func (cmd *Hunt) Run(cli *cli.Globals) error {
	if cli.Help {
		fmt.Print(Usage)
		return nil
	}

	var schema stream.Streamable

	var db *hunt.Database
	var fn text.Colored

	cli.NoConvert = true // force

	hs := cli.Bootstrap(cmd.Paths)
	defer cli.ThrowAway()

	if cli.Verbose > 0 {
		log.Println("hunt: started")
	}

	if cmd.Sqlite {
		db = hunt.NewDB(types.Database)

		if cli.Verbose > 0 {
			log.Printf("hunt: using database %s\n", db)
		}
	}

	if len(cmd.Url) > 0 {
		switch {
		case cmd.Hec:
			schema = hec.New(cmd.Url, cmd.Auth)
		case cmd.Ecs:
			schema = ecs.New(cmd.Url)
		default:
			schema = raw.New(cmd.Url)
		}

		if cli.Verbose > 0 {
			log.Printf("hunt: using schema %s\n", schema)
		}
	}

	cnt := 0

	for e := range hunt.Hunt(hs, &hunt.Options{
		Sort:       cmd.Sort,
		Extensions: cmd.Ext,
		Verbose:    cli.Verbose,
	}) {
		if cmd.All || e.Severity >= hunt.Level {
			switch {
			case cmd.All && e.Severity >= hunt.Level:
				fn = text.Mark // mark event
			case cmd.All:
				fn = text.Hide // hide event
			default:
				fn = text.Term // reset terminal
			}

			switch {
			case cmd.Jsonl:
				_, _ = fmt.Fprintln(cli.Stdout, fn(e.ToJSONL()))
			case cmd.Json:
				_, _ = fmt.Fprintln(cli.Stdout, fn(e.ToJSON()))
			default:
				_, _ = fmt.Fprintln(cli.Stdout, fn(e.ToCEF()))
			}

			// TODO: concurrent
			if db != nil {
				db.Write(e)
			}

			// TODO: concurrent
			if len(cmd.Url) > 0 {
				_ = schema.Write(e)
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
