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
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/formats"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/schemas"
	"go.foxforensics.eu/fox/v4/internal/pkg/types"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/client"
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
  -U, --url=URL            Stream events with CEF schema to URL
  -E, --ecs=URL            Stream events with ECS schema to URL
  -H, --hec=URL            Stream events with HEC schema to URL
  -A, --auth=TOKEN         Use auth token for HEC streaming

Remarks:
  If 'local' is specified as path, built-in paths will be used.

Example: Hunt down critical events
  $ fox hunt -u *.dd

Example: Save local events as Parquet
  $ fox hunt -ap local

Example: Send events to an Elastic Stack
  $ fox hunt -E http://127.0.0.1:8080 *.evtx

Report bugs at: foxforensics.eu/issues
`)

//go:embed hunt.yml
var Default []byte

type Hunt struct {
	All     bool `short:"a"`
	Sort    bool `short:"s"`
	Uniq    bool `short:"u" xor:"uniq"`
	Json    bool `short:"j" xor:"json,jsonl"`
	Jsonl   bool `short:"J" xor:"json,jsonl"`
	Parquet bool `short:"p"`

	// sigma flags
	Rule []byte `short:"R" type:"filecontent"`

	// stream flags
	Url  string `short:"U" xor:"url,ecs,hec"`
	Ecs  string `short:"E" xor:"url,ecs,hec"`
	Hec  string `short:"H" xor:"url,ecs,hec"`
	Auth string `short:"A" and:"auth,hec"`

	// paths
	Paths []string `arg:"" optional:""`

	// internal
	net  *client.Client   `kong:"-"`
	file *parquet.Parquet `kong:"-"`
	uniq *types.Unique    `kong:"-"`
	rule sigma.Rule       `kong:"-"`
}

func (cmd *Hunt) Validate() error {
	if len(cmd.Hec) > 0 && len(cmd.Auth) == 0 {
		return errors.New("auth required")
	}

	if len(cmd.Hec) == 0 && len(cmd.Auth) > 0 {
		return errors.New("no auth required")
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

	if len(cmd.Rule) > 0 {
		cmd.rule, err = sigma.ParseRule(cmd.Rule)
	} else {
		cmd.rule, err = sigma.ParseRule(Default)
	}

	if err != nil {
		return err
	}

	if len(cmd.Url) > 0 {
		cmd.net, err = client.New(&client.Options{
			Url:    cmd.Url,
			Schema: schemas.Raw,
		})
	}

	if len(cmd.Ecs) > 0 {
		cmd.net, err = client.New(&client.Options{
			Url:    cmd.Ecs,
			Schema: schemas.Ecs,
		})
	}

	if len(cmd.Hec) > 0 {
		cmd.net, err = client.New(&client.Options{
			Url:    cmd.Hec,
			Token:  cmd.Auth,
			Schema: schemas.Hec,
		})
	}

	return err
}

func (cmd *Hunt) Run(fox *cmd.Globals) error {
	cmd.Paths = append(cmd.Paths, fox.Paths...)

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

	if cmd.net != nil {
		slog.Debug(fmt.Sprintf("hunt: streaming to %s", cmd.net))
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
			fox.Writer.Match(formats.Event(e, cmd.Json, cmd.Jsonl), fox.Regexp)
		} else {
			if err = cmd.file.Write(e); err != nil {
				slog.Error(err.Error())
			}
		}

		if cmd.net != nil {
			if err = cmd.net.Stream(fox.Context, e); err != nil {
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
