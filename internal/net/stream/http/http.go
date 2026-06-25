package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/cenkalti/backoff"
	"go.foxforensics.eu/fox/v4/internal/net/client"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/schemas"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/schemas/ecs"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/schemas/hec"
	"go.foxforensics.eu/fox/v4/internal/pkg/adapters/schemas/raw"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/event"
)

type Options struct {
	Url    string
	Token  string
	Schema schemas.Schema
}

type Http struct {
	client *http.Client
	opts   *Options
}

func Create(opts *Options) (*Http, error) {
	u, err := url.Parse(opts.Url)

	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(u.Scheme, "http") {
		return nil, errors.New("unsupported scheme")
	}

	if u.Scheme == "http" {
		slog.Warn("data will be streamed unencrypted!")
	}

	return &Http{client.Http(), opts}, nil
}

func (h Http) String() string {
	return fmt.Sprintf("%s@%s", h.opts.Schema, h.opts.Url)
}

func (h Http) Stream(ctx context.Context, evt *event.Event) error {
	var buf []byte
	var err error

	switch h.opts.Schema {
	case schemas.Ecs:
		buf, err = ecs.Apply(evt)
	case schemas.Hec:
		buf, err = hec.Apply(evt)
	case schemas.Raw:
		buf, err = raw.Apply(evt)
	}

	if err != nil {
		return err
	}

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = client.Timeout

	return backoff.Retry(func() error {
		req, err := http.NewRequestWithContext(ctx, "POST", h.opts.Url, bytes.NewReader(buf))

		if err != nil {
			return err
		}

		req.Header.Add("User-Agent", client.Name())

		switch h.opts.Schema {
		case schemas.Ecs, schemas.Hec:
			req.Header.Set("Content-Type", "application/json")
		default:
			req.Header.Set("Content-Type", "text/plain")
		}

		// add authorization for Splunk
		if h.opts.Schema == schemas.Hec && len(h.opts.Token) > 0 {
			req.Header.Set("Authorization", fmt.Sprintf("Splunk %s", h.opts.Token))
		}

		res, err := h.client.Do(req)

		if err != nil {
			return err
		}

		defer func() {
			_ = res.Body.Close()
		}()

		// drain body
		_, _ = io.Copy(io.Discard, res.Body)

		switch {
		case res.StatusCode >= 500: // retry
			return errors.New(res.Status)
		case res.StatusCode >= 400: // halt
			return backoff.Permanent(errors.New(res.Status))
		default: // success
			return nil
		}
	}, bo)
}

func (h Http) Close() error {
	return nil // stateless
}
