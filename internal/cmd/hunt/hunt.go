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
	"go.foxforensics.eu/fox/v5/internal/cmd"
	"go.foxforensics.eu/fox/v5/internal/cmd/hunt/client"
	"go.foxforensics.eu/fox/v5/internal/cmd/hunt/event"
	"go.foxforensics.eu/fox/v5/internal/cmd/hunt/hunter"
	"go.foxforensics.eu/fox/v5/internal/cmd/hunt/muxer"
	"go.foxforensics.eu/fox/v5/internal/cmd/hunt/parquet"
	"go.foxforensics.eu/fox/v5/internal/cmd/hunt/rules"
	"go.foxforensics.eu/fox/v5/internal/pkg"
	"go.foxforensics.eu/fox/v5/internal/pkg/receipt"
	"go.foxforensics.eu/fox/v5/internal/pkg/types"
	"go.foxforensics.eu/fox/v5/library/formats"
)

var Usage = strings.TrimSpace(`
Usage: fox hunt [FLAGS...] <local|PATHS...>

Flags:
  -a, --all                Show logs with all severities
  -s, --sort               Show logs sorted by timestamp
  -u, --uniq               Show logs that are unique 
  -j, --json               Show logs as JSON objects
  -l, --jsonl              Show logs as JSON lines
  -t, --triage             Show logs as Triage
  -p, --parquet            Save logs as Parquet

Filter flags:
  -N, --min=TIME           Minimum event time (RFC3339)
  -X, --max=TIME           Maximum event time (RFC3339)
  -R, --rule=FILE          Filter using Sigma rules file

Stream flags:
  -U, --url=URL            Stream events using CEF schema
  -E, --ecs=URL            Stream events using ECS schema
  -H, --hec=URL            Stream events using HEC schema
  -A, --auth=TOKEN         Use auth token with HEC streaming

Remarks:
  If 'local' is specified as path, built-in paths will be used.

Example: Hunt down critical events
  $ fox hunt -t *.dd

Example: Save local events as Parquet
  $ fox hunt -ap local

Example: Send events to an Elastic Stack
  $ fox hunt -E http://127.0.0.1:8080 *.evtx

Report bugs at: foxforensics.eu/issues
`)

type Hunt struct {
	All     bool `short:"a"`
	Sort    bool `short:"s"`
	Uniq    bool `short:"u"`
	Json    bool `short:"j" xor:"triage,json,jsonl"`
	Jsonl   bool `short:"l" xor:"triage,json,jsonl"`
	Triage  bool `short:"t" xor:"triage,json,jsonl"`
	Parquet bool `short:"p"`

	// filter flags
	Min  time.Time `short:"N"`
	Max  time.Time `short:"X"`
	Rule []byte    `short:"R" type:"filecontent"`

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
	if !cmd.Min.IsZero() && !cmd.Max.IsZero() && cmd.Min.After(cmd.Max) {
		return errors.New("invalid range")
	}

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

	if cmd.Triage {
		cmd.Sort = true
		cmd.Uniq = true
	}

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
		cmd.rule, err = sigma.ParseRule(rules.Critical)
	}

	if err != nil {
		return err
	}

	switch {
	case len(cmd.Url) > 0:
		cmd.client, err = client.Raw(cmd.Url)
	case len(cmd.Ecs) > 0:
		cmd.client, err = client.Ecs(cmd.Ecs)
	case len(cmd.Hec) > 0:
		cmd.client, err = client.Hec(cmd.Hec, cmd.Auth)
	}

	return err
}

func (cmd *Hunt) Run(fox *cmd.Globals) error {
	cmd.Paths = append(cmd.Paths, fox.Paths...)

	if len(cmd.Paths) == 0 {
		pkg.Usage(Usage)
		return nil
	}

	if cmd.Paths[0] == "local" {
		cmd.Paths = append(hunter.Local, cmd.Paths[1:]...)
	}

	heaps, err := fox.Init(cmd.Paths, true)

	if err != nil {
		return err
	}

	defer cmd.discard(fox)

	sig := evaluator.ForRule(cmd.rule)
	mux := muxer.New(fox.Context, hunter.Scale)

	slog.Info("hunt: started")
	slog.Debug(fmt.Sprintf("hunt: using %d thread(s)", fox.Threads))
	slog.Debug(fmt.Sprintf("hunt: using rule '%s'", cmd.rule.Title))

	if cmd.rule.Logsource.Product != "fox" {
		slog.Warn("hunt: rule is not officially supported!")
	}

	if cmd.client != nil {
		slog.Debug(fmt.Sprintf("hunt: streaming to %s", cmd.client))
		mux.AddHandler(cmd.client.Send)
	}

	if cmd.parquet != nil {
		slog.Debug(fmt.Sprintf("hunt: using storage %s", cmd.parquet))
		mux.AddHandler(cmd.parquet.Write)
	}

	if !fox.Quiet {
		mux.AddHandler(func(_ context.Context, e *event.Event) error {
			fox.Writer.Match(formats.Event(e, cmd.Json, cmd.Jsonl, cmd.Triage), fox.Regexp)
			return nil
		})
	}

	var n int64

	for e := range hunter.New(&hunter.Options{
		Sort:    cmd.Sort,
		Threads: fox.Threads,
	}).Hunt(fox.Context, heaps) {
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

		if !cmd.Min.IsZero() && e.Time.Before(cmd.Min) {
			continue // to soon
		}

		if !cmd.Max.IsZero() && e.Time.After(cmd.Max) {
			continue // to late
		}

		mux.Handle(fox.Context, e)
		n++
	}

	slog.Info("hunt: finished")
	slog.Info(fmt.Sprintf("hunt: found %d event(s)", n))

	return mux.Wait()
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
