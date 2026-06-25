package hunt

import (
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/bradleyjkemp/sigma-go"
	"github.com/bradleyjkemp/sigma-go/evaluator"
	"go.foxforensics.eu/fox/v4/internal/cmd"
	"go.foxforensics.eu/fox/v4/internal/net/stream"
	"go.foxforensics.eu/fox/v4/internal/net/stream/http"
	"go.foxforensics.eu/fox/v4/internal/net/stream/mqtt"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/formats"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/schemas"
	"go.foxforensics.eu/fox/v4/internal/pkg/types"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/hunter"
	"go.foxforensics.eu/fox/v4/internal/sys"
	"go.foxforensics.eu/fox/v4/internal/sys/parquet"
	"go.foxforensics.eu/fox/v4/internal/sys/receipt"
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

Sigma flags:
  -R, --rule=FILE          Filter using Sigma rules file

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

//go:embed hunt.yml
var Default []byte

type Hunt struct {
	All     bool `short:"a"`
	Sort    bool `short:"s"`
	Uniq    bool `short:"u" xor:"uniq,dist"`
	Json    bool `short:"j" xor:"json,jsonl"`
	Jsonl   bool `short:"J" xor:"json,jsonl"`
	Parquet bool `short:"p"`

	// sigma flags
	Rule []byte `short:"R" type:"filecontent"`

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
	streamer stream.Streamer  `kong:"-"`
	file     *parquet.Parquet `kong:"-"`
	uniq     *types.Unique    `kong:"-"`
	rule     sigma.Rule       `kong:"-"`
}

func (cmd *Hunt) Validate() error {
	if cmd.Hec && len(cmd.Auth) == 0 && len(cmd.Mqtt) == 0 {
		return errors.New("auth required")
	}

	if cmd.Hec && len(cmd.Auth) > 0 && len(cmd.Mqtt) > 0 {
		return errors.New("no auth required")
	}

	if cmd.QoS > 2 {
		return errors.New("mqtt qos invalid")
	}

	return nil
}

func (cmd *Hunt) AfterApply(_ *kong.Kong, _ kong.Vars) error {
	var err error

	if cmd.Uniq {
		cmd.uniq = types.NewUnique()
	}

	if cmd.Parquet {
		cmd.file, err = parquet.New(fmt.Sprintf("fox_hunt_%s",
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
		var shm schemas.Schema

		switch {
		case cmd.Hec:
			shm = schemas.Hec
		case cmd.Ecs:
			shm = schemas.Ecs
		default:
			shm = schemas.Raw
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

	if len(cmd.Rule) > 0 {
		cmd.rule, err = sigma.ParseRule(cmd.Rule)
	} else {
		cmd.rule, err = sigma.ParseRule(Default)
	}

	return err
}

func (cmd *Hunt) Run(fox *cmd.Globals) error {
	cmd.Paths = append(cmd.Paths, fox.Input...)

	if len(cmd.Paths) == 0 {
		return sys.Usage(Usage)
	}

	if cmd.Paths[0] == "local" {
		cmd.Paths = append(hunter.Local, cmd.Paths[1:]...)
	}

	ch, err := fox.Init(cmd.Paths, true)

	if err != nil {
		return err
	}

	defer cmd.discard(fox)

	slog.Info("hunt: started")
	slog.Debug(fmt.Sprintf("hunt: using %d thread(s)", fox.Threads))
	slog.Debug(fmt.Sprintf("hunt: using rule '%s'", cmd.rule.Title))

	if cmd.rule.Logsource.Product != "fox" {
		slog.Warn("hunt: rule is not officially supported!")
	}

	if cmd.file != nil {
		slog.Debug(fmt.Sprintf("hunt: using storage %s", cmd.file))
	}

	if cmd.streamer != nil {
		slog.Debug(fmt.Sprintf("hunt: streaming as %s", cmd.streamer))
	}

	var n int64

	var sig = evaluator.ForRule(cmd.rule)

	for e := range hunter.New(&hunter.Options{
		Sort:    cmd.Sort,
		Threads: fox.Threads,
	}).Hunt(fox.Context, ch) {
		m, err := sig.Matches(fox.Context, e.Fields)

		if err != nil {
			slog.Error(err.Error())
			continue // not successful
		}

		if cmd.uniq != nil && !cmd.uniq.Is(e.String()) {
			continue // not unique
		}

		if !cmd.All && !m.Match {
			continue // not matched
		}

		if cmd.file == nil {
			fox.Stdout.Match(formats.Event(e, cmd.Json, cmd.Jsonl), fox.Regexp)
		} else {
			if err = cmd.file.Write(e); err != nil {
				slog.Error(err.Error())
			}
		}

		if cmd.streamer != nil {
			err = cmd.streamer.Stream(fox.Context, e)

			if err != nil {
				slog.Error(err.Error())
			}
		}

		n++
	}

	slog.Info("hunt: finished")
	slog.Info(fmt.Sprintf("hunt: found %d event(s)", n))

	return nil
}

func (cmd *Hunt) discard(fox *cmd.Globals) {
	if cmd.streamer != nil {
		err := cmd.streamer.Close()

		if err != nil {
			slog.Error(err.Error())
		}
	}

	if cmd.file != nil {
		err := cmd.file.Close()

		if err != nil {
			slog.Error(err.Error())
		}

		if fox.NoReceipt {
			return
		}

		err = receipt.Generate(cmd.file.String())

		if err != nil {
			slog.Error(err.Error())
		}
	}
}
