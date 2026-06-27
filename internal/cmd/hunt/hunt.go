package hunt

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/bradleyjkemp/sigma-go"
	"github.com/bradleyjkemp/sigma-go/evaluator"
	"github.com/sourcegraph/conc/pool"
	"go.foxforensics.eu/fox/v4/internal/cmd"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/formats"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/schemas"
	"go.foxforensics.eu/fox/v4/internal/pkg/types"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/client"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/event"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/hunter"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/parquet"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/receipt"
	"go.foxforensics.eu/fox/v4/internal/sys"
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
  -U, --url=URL            Stream events using CEF schema
  -E, --ecs=URL            Stream events using ECS schema
  -H, --hec=URL            Stream events using HEC schema
  -A, --auth=TOKEN         Use auth token with HEC streaming

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
	Uniq    bool `short:"u"`
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
	client  *client.Client   `kong:"-"`
	parquet *parquet.Parquet `kong:"-"`
	unique  *types.Unique    `kong:"-"`
	rule    sigma.Rule       `kong:"-"`
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
		cmd.unique = types.NewUnique()
	}

	if cmd.Parquet {
		cmd.parquet, err = parquet.New(fmt.Sprintf("fox_hunt_%s",
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

	switch {
	case len(cmd.Url) > 0:
		cmd.client, err = client.New(&client.Options{
			Url:    cmd.Url,
			Schema: schemas.Raw,
		})

	case len(cmd.Ecs) > 0:
		cmd.client, err = client.New(&client.Options{
			Url:    cmd.Ecs,
			Schema: schemas.Ecs,
		})

	case len(cmd.Hec) > 0:
		cmd.client, err = client.New(&client.Options{
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

	p := pool.New().
		WithContext(fox.Context).
		WithMaxGoroutines(fox.Threads)

	ch1 := make(chan *event.Event, fox.Threads*hunter.Scale)
	ch2 := make(chan *event.Event, fox.Threads*hunter.Scale)

	slog.Info("hunt: started")
	slog.Debug(fmt.Sprintf("hunt: using %d thread(s)", fox.Threads))
	slog.Debug(fmt.Sprintf("hunt: using rule '%s'", cmd.rule.Title))

	if cmd.rule.Logsource.Product != "fox" {
		slog.Warn("hunt: rule is not officially supported!")
	}

	if cmd.parquet != nil {
		slog.Debug(fmt.Sprintf("hunt: using storage %s", cmd.parquet))

		p.Go(func(ctx context.Context) error {
			return cmd.parquet.Run(ctx, ch1)
		})
	}

	if cmd.client != nil {
		slog.Debug(fmt.Sprintf("hunt: streaming to %s", cmd.client))

		p.Go(func(ctx context.Context) error {
			return cmd.client.Run(ctx, ch2)
		})
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

		if cmd.unique != nil && !cmd.unique.Is(e.String()) {
			continue // not unique
		}

		if !cmd.All && !m.Match {
			continue // not matched
		}

		if cmd.parquet != nil {
			ch1 <- e
		}

		if cmd.client != nil {
			ch2 <- e
		}

		fox.Writer.Match(formats.Event(e, cmd.Json, cmd.Jsonl), fox.Regexp)

		n++
	}

	close(ch1)
	close(ch2)

	slog.Info("hunt: finished")
	slog.Info(fmt.Sprintf("hunt: found %d event(s)", n))

	return p.Wait()
}

func (cmd *Hunt) discard(fox *cmd.Globals) {
	if cmd.parquet == nil {
		return
	}

	err := cmd.parquet.Close()

	if err != nil {
		slog.Error(err.Error())
	}

	if fox.NoReceipt {
		return
	}

	err = receipt.Generate(cmd.parquet.String())

	if err != nil {
		slog.Error(err.Error())
	}
}
