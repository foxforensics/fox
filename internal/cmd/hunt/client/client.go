package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v6"
	"go.foxforensics.eu/fox/v5/internal/cmd/hunt/client/ecs"
	"go.foxforensics.eu/fox/v5/internal/cmd/hunt/client/hec"
	"go.foxforensics.eu/fox/v5/internal/cmd/hunt/client/raw"
	"go.foxforensics.eu/fox/v5/internal/cmd/hunt/event"
	"go.foxforensics.eu/fox/v5/internal/pkg/version"
)

// Timeout for everything network related.
const timeout = time.Second * 30

type formater func(*event.Event) ([]byte, error)

type Format int

const (
	Raw Format = iota
	ECS
	HEC
)

func (f Format) String() string {
	switch f {
	case Raw:
		return "raw"
	case ECS:
		return "ecs"
	case HEC:
		return "hec"
	default:
		return "unknown"
	}
}

type Options struct {
	Url    string
	Token  string
	Format Format
}

type Client struct {
	client http.Client
	header http.Header
	apply  formater
	opts   *Options
}

func New(opts *Options) (*Client, error) {
	u, err := url.Parse(opts.Url)

	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(u.Scheme, "http") {
		return nil, errors.New("scheme not supported")
	}

	if u.Scheme == "http" {
		slog.Warn("events will be streamed unencrypted!")
	}

	cli := &Client{
		client: http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				Proxy:               http.ProxyFromEnvironment,
				IdleConnTimeout:     timeout,
				TLSHandshakeTimeout: timeout,
				MaxIdleConnsPerHost: 10,
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS13, // pinned
				},
			},
		},
		header: make(http.Header, 4),
		opts:   opts,
	}

	// add fox agent
	cli.header.Set("User-Agent", fmt.Sprintf("fox %s", version.Number))

	// add Splunk token
	if len(cli.opts.Token) > 0 {
		cli.header.Set("Authorization", fmt.Sprintf("Splunk %s", cli.opts.Token))
	}

	// add content type
	if cli.opts.Format != Raw {
		cli.header.Set("Content-Type", "application/json")
	} else {
		cli.header.Set("Content-Type", "text/plain")
	}

	switch cli.opts.Format {
	case ECS:
		cli.apply = ecs.Apply
	case HEC:
		cli.apply = hec.Apply
	case Raw:
		cli.apply = raw.Apply
	}

	return cli, nil
}

func WithRaw(url string) (*Client, error) {
	return New(&Options{
		Url:    url,
		Format: Raw,
	})
}

func WithEcs(url string) (*Client, error) {
	return New(&Options{
		Url:    url,
		Format: ECS,
	})
}

func WithHec(url, token string) (*Client, error) {
	return New(&Options{
		Url:    url,
		Token:  token,
		Format: HEC,
	})
}

func (cli *Client) String() string {
	return fmt.Sprintf("%s@%s", cli.opts.Format, cli.opts.Url)
}

func (cli *Client) Send(ctx context.Context, evt *event.Event) error {
	buf, err := cli.apply(evt)

	if err != nil {
		return err
	}

	_, err = backoff.Retry(ctx, func() (any, error) {
		req, err := http.NewRequestWithContext(ctx, "POST", cli.opts.Url, bytes.NewReader(buf))

		if err != nil {
			return nil, err
		}

		req.Header = cli.header

		res, err := cli.client.Do(req)

		if err != nil {
			return nil, err
		}

		defer func() {
			if err := res.Body.Close(); err != nil {
				slog.Error(err.Error())
			}
		}()

		// drain body
		_, err = io.Copy(io.Discard, res.Body)

		if err != nil {
			slog.Error(err.Error())
		}

		switch {
		case res.StatusCode >= 500: // retry
			return nil, errors.New(res.Status)
		case res.StatusCode >= 400: // halt
			return nil, backoff.Permanent(errors.New(res.Status))
		default:
			return nil, nil
		}
	},
		backoff.WithMaxElapsedTime(timeout),
		backoff.WithMaxTries(10),
	)

	if err, ok := errors.AsType[*backoff.RetryError](err); ok {
		return err.LastErr // unwrap client error
	}

	return err
}
