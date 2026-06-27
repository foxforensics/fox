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
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/schemas"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/schemas/ecs"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/schemas/hec"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/schemas/raw"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/event"
	"go.foxforensics.eu/fox/v4/internal/sys/version"
)

// Timeout for everything network related.
const Timeout = time.Second * 30

type Options struct {
	Url    string
	Token  string
	Schema schemas.Schema
}

type Client struct {
	client http.Client
	header http.Header
	apply  schemas.Apply
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
			Timeout: Timeout,
			Transport: &http.Transport{
				Proxy:               http.ProxyFromEnvironment,
				IdleConnTimeout:     Timeout,
				TLSHandshakeTimeout: Timeout,
				MaxIdleConnsPerHost: 10,
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS13, // pinned
				},
			},
		},
		header: make(http.Header, 4),
		opts:   opts}

	// add fox agent
	cli.header.Set("User-Agent", fmt.Sprintf("fox %s", version.Number))

	// add Splunk token
	if len(cli.opts.Token) > 0 {
		cli.header.Set("Authorization", fmt.Sprintf("Splunk %s", cli.opts.Token))
	}

	// add content type
	if cli.opts.Schema != schemas.Raw {
		cli.header.Set("Content-Type", "application/json")
	} else {
		cli.header.Set("Content-Type", "text/plain")
	}

	switch cli.opts.Schema {
	case schemas.Ecs:
		cli.apply = ecs.Apply
	case schemas.Hec:
		cli.apply = hec.Apply
	case schemas.Raw:
		cli.apply = raw.Apply
	}

	return cli, nil
}

func (cli *Client) Run(ctx context.Context, ch <-chan *event.Event) error {
	for e := range ch {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			if err := cli.stream(ctx, e); err != nil {
				slog.Error(err.Error())
			}
		}
	}

	return nil
}

func (cli *Client) String() string {
	return fmt.Sprintf("%s@%s", cli.opts.Schema, cli.opts.Url)
}

func (cli *Client) stream(ctx context.Context, evt *event.Event) error {
	buf, err := cli.apply(evt)

	if err != nil {
		return err
	}

	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = Timeout

	_, err = backoff.Retry(ctx, func() (any, error) {
		select {
		case <-ctx.Done():
			return nil, backoff.Permanent(ctx.Err())
		default:
		}

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
	}, backoff.WithBackOff(bo))

	return err
}
