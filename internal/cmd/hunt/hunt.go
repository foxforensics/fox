package hunt

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/bradleyjkemp/sigma-go"
	"github.com/bradleyjkemp/sigma-go/evaluator"

	cli "go.foxforensics.eu/fox/v4/internal/cmd"

	"go.foxforensics.eu/fox/v4/internal/pkg/file/schema"
	"go.foxforensics.eu/fox/v4/internal/pkg/file/storage"
	"go.foxforensics.eu/fox/v4/internal/pkg/file/storage/parquet"
	"go.foxforensics.eu/fox/v4/internal/pkg/file/stream"
	"go.foxforensics.eu/fox/v4/internal/pkg/file/stream/http"
	"go.foxforensics.eu/fox/v4/internal/pkg/file/stream/mqtt"
	"go.foxforensics.eu/fox/v4/internal/pkg/rules"
	"go.foxforensics.eu/fox/v4/internal/pkg/text"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/event"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/hunter"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/receipt"
)

var Usage = strings.TrimSpace(`
Usage: fox hunt [FLAGS...] <local|PATHS...>

Flags:
  -a, --all                Show logs with all severities
  -s, --sort               Show logs sorted by timestamp (slow)
  -u, --uniq               Show logs that are unique 
  -j, --json               Show logs as JSON objects
  -J, --jsonl              Show logs as JSON lines
  -p, --parquet            Save logs as Parquet

Filter flags:
  -R, --rule=FILE          Filter using Sigma Rules file
  -D, --dist=LENGTH        Filter using Levenshtein distance (slow)

Stream flags:
  -U, --url=SERVER         Stream events to a server or broker
  -A, --auth=TOKEN         Use token for streaming to Splunk
  -M, --mqtt=TOPIC         Use topic for streaming via MQTT

Schema flags:
  -e, --ecs                Use ECS schema while streaming
  -h, --hec                Use HEC schema while streaming

Aliases:
  --elastic                Alias for -eU http://localhost:8080
  --splunk                 Alias for -hU http://localhost:8088/...

Remarks:
  If 'local' is specified as path, built-in paths will be used.

Example: Hunt down critical events
  $ fox hunt -u *.dd

Example: Save all events as Parquet
  $ fox hunt -ap *.evtx

Example: Send local events to a server
  $ fox hunt -U http://127.0.0.1:8080 local

Example: Send local events to a broker
  $ fox hunt -U tcp://127.0.0.1:1883 -M events local

Report bugs at: foxforensics.eu/issues
`)

type Hunt struct {
	All     bool `short:"a"`
	Sort    bool `short:"s"`
	Uniq    bool `short:"u" xor:"uniq,dist"`
	Json    bool `short:"j" xor:"json,jsonl"`
	Jsonl   bool `short:"J" xor:"json,jsonl"`
	Parquet bool `short:"p"`

	// filter flags
	Rule string  `short:"R"`
	Dist float64 `short:"D" xor:"uniq,dist"`

	// stream flags
	Url  string `short:"U"`
	Auth string `short:"A" xor:"auth,mqtt"`
	Mqtt string `short:"M" xor:"auth,mqtt"`

	// schema flags
	Ecs bool `short:"e" xor:"ecs,hec"`
	Hec bool `short:"h" xor:"ecs,hec"`

	// aliases
	Elastic bool `long:"elastic" xor:"elastic,splunk"`
	Splunk  bool `long:"splunk" xor:"elastic,splunk"`

	// hidden
	QoS      byte   `hidden:"" long:"mqtt-qos"`
	Username string `hidden:"" long:"mqtt-username"`
	Password string `hidden:"" long:"mqtt-password"`

	// paths
	Paths []string `arg:"" optional:""`

	// internal
	storage  storage.Storage `kong:"-"`
	streamer stream.Streamer `kong:"-"`
	rule     sigma.Rule      `kong:"-"`
	uniq     text.Unique     `kong:"-"`
}

func (cmd *Hunt) Validate() error {
	if cmd.Hec && (len(cmd.Auth)+len(cmd.Mqtt) == 0) {
		return errors.New("auth required")
	}

	return nil
}

func (cmd *Hunt) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	var err error

	switch {
	case cmd.Uniq:
		cmd.uniq = text.ByHash()
	case cmd.Dist > 0:
		cmd.uniq = text.ByDistance(cmd.Dist)
	}

	if cmd.Parquet {
		cmd.storage, err = parquet.Create(fmt.Sprintf("fox_hunt_%s",
			time.Now().UTC().Format("20060102150405"),
		))

		if err != nil {
			return err
		}
	}

	if cmd.Elastic {
		cmd.Url = stream.Elastic
		cmd.Ecs = true
	}

	if cmd.Splunk {
		cmd.Url = stream.Splunk
		cmd.Hec = true
	}

	if len(cmd.Url) > 0 {
		var shm schema.Schema

		switch {
		case cmd.Hec:
			shm = schema.Hec
		case cmd.Ecs:
			shm = schema.Ecs
		default:
			shm = schema.Raw
		}

		if len(cmd.Mqtt) > 0 {
			cmd.streamer, err = mqtt.Create(&mqtt.Options{
				Url:      cmd.Url,
				Topic:    cmd.Mqtt,
				Username: cmd.Username,
				Password: cmd.Password,
				QoS:      cmd.QoS,
				Schema:   shm,
			})
		} else {
			cmd.streamer, err = http.Create(&http.Options{
				Url:    cmd.Url,
				Token:  cmd.Auth,
				Schema: shm,
			})
		}

		if err != nil {
			return err
		}
	}

	rule := rules.Critical

	if len(cmd.Rule) > 0 {
		rule, err = os.ReadFile(cmd.Rule)

		if err != nil {
			return err
		}
	}

	cmd.rule, err = sigma.ParseRule(rule)

	return err
}

func (cmd *Hunt) Run(cli *cli.Globals) error {
	cmd.Paths = append(cmd.Paths, cli.Input...)

	if len(cmd.Paths) == 0 {
		return text.Usage(Usage)
	}

	if cmd.Paths[0] == "local" {
		cmd.Paths = append(hunter.Local, cmd.Paths[1:]...)
	}

	ch, err := cli.Init(cmd.Paths, true)

	if err != nil {
		return err
	}

	defer cli.Discard()
	defer cmd.discard(cli)

	if cli.Verbose > 0 {
		slog.Info("hunt: started")
	}

	if cli.Verbose > 1 {
		slog.Info(fmt.Sprintf("hunt: using %d thread(s)", cli.Threads))
	}

	if cli.Verbose > 1 {
		slog.Info(fmt.Sprintf("hunt: using rule '%s'", cmd.rule.Title))
	}

	if cli.Verbose > 1 && cmd.storage != nil {
		slog.Info(fmt.Sprintf("hunt: using storage %s", cmd.storage))
	}

	if cli.Verbose > 1 && cmd.streamer != nil {
		slog.Info(fmt.Sprintf("hunt: streaming as %s", cmd.streamer))
	}

	if !rules.IsSupported(&cmd.rule) {
		slog.Warn("rule is not supported")
	}

	var n int64

	var sig = evaluator.ForRule(cmd.rule)

	for e := range hunter.New(&hunter.Options{
		Sort:    cmd.Sort,
		Threads: cli.Threads,
		Verbose: cli.Verbose,
	}).Hunt(cli.Context, ch) {
		m, err := sig.Matches(cli.Context, e.Fields)

		if err != nil {
			slog.Error(err.Error())
			continue // not successful
		}

		if cmd.uniq != nil && !cmd.uniq.IsUnique(e.String()) {
			continue // not unique
		}

		if !cmd.All && !m.Match {
			continue // not matched
		}

		if cmd.storage == nil {
			text.Stdout.Match(cmd.format(e), cli.Regexp)
		} else {
			err = cmd.storage.Store(e)

			if err != nil {
				slog.Error(err.Error())
			}
		}

		if cmd.streamer != nil {
			err = cmd.streamer.Stream(e)

			if err != nil {
				slog.Error(err.Error())
			}
		}

		n++
	}

	if cli.Verbose > 0 {
		slog.Info("hunt: finished")
	}

	if cli.Verbose > 1 {
		slog.Info(fmt.Sprintf("hunt: found %d event(s)", n))
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
	if cmd.streamer != nil {
		err := cmd.streamer.Close()

		if err != nil {
			slog.Error(err.Error())
		}
	}

	if cmd.storage != nil {
		err := cmd.storage.Close()

		if err != nil {
			slog.Error(err.Error())
		}

		if cli.NoReceipt {
			return
		}

		err = receipt.Generate(cmd.storage.String())

		if err != nil {
			slog.Error(err.Error())
		}
	}
}
